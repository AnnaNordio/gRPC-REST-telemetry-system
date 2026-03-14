package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"
	"github.com/anna/iot-dual-stack/gen/sensor" // Assicurati che il path sia corretto
)

// Struttura per memorizzare i dati ricevuti (per statistiche tesi)
type Stats struct {
	mu           sync.Mutex
	RestCount    int
	GrpcCount    int
	TotalPayload int
}

var globalStats Stats

// --- IMPLEMENTAZIONE gRPC ---
type backendGrpcServer struct {
	sensor.UnimplementedTelemetryServiceServer
}

func (s *backendGrpcServer) SendData(ctx context.Context, in *sensor.TelemetryData) (*sensor.Reply, error) {
	globalStats.mu.Lock()
	globalStats.GrpcCount++
	globalStats.mu.Unlock()

	log.Printf("[BE-gRPC] Ricevuto dato dal sensore: %s | Temp: %.2f", in.SensorId, in.Temperature)
	
	return &sensor.Reply{
		Status:     "SUCCESS_GRPC",
		ServerTime: time.Now().Unix(),
	}, nil
}

// --- IMPLEMENTAZIONE REST ---
func restBackendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Solo POST ammesso", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Errore lettura", http.StatusInternalServerError)
		return
	}

	var data sensor.TelemetryData
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("[BE-REST] Errore decoding: %v", err)
		return
	}

	globalStats.mu.Lock()
	globalStats.RestCount++
	globalStats.TotalPayload += len(body)
	globalStats.mu.Unlock()

	log.Printf("[BE-REST] Ricevuto dato dal sensore: %s | Dimensione: %d bytes", data.SensorId, len(body))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "SUCCESS_REST"})
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	// 1. Avvio Server gRPC (Porta 50052 - Nota: diversa dal Gateway)
	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", ":50052")
		if err != nil {
			log.Fatalf("Errore listener gRPC: %v", err)
		}
		
		s := grpc.NewServer()
		sensor.RegisterTelemetryServiceServer(s, &backendGrpcServer{})
		
		log.Println("✅ Backend gRPC in ascolto su :50052")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Errore server gRPC: %v", err)
		}
	}()

	// 2. Avvio Server REST (Porta 8081 - Nota: diversa dal Gateway)
	go func() {
		defer wg.Done()
		mux := http.NewServeMux()
		mux.HandleFunc("/process", restBackendHandler)
		
		log.Println("✅ Backend REST in ascolto su :8081")
		if err := http.ListenAndServe(":8081", mux); err != nil {
			log.Fatalf("Errore server REST: %v", err)
		}
	}()

	// Endpoint opzionale per vedere le statistiche accumulate
	go func() {
		http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
			globalStats.mu.Lock()
			json.NewEncoder(w).Encode(globalStats)
			globalStats.mu.Unlock()
		})
		http.ListenAndServe(":9090", nil)
	}()

	wg.Wait()
}

