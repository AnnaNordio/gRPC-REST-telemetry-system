package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"github.com/anna/iot-dual-stack/gen/sensor" // Assicurati che il path sia corretto
)

// Server gRPC
type grpcServer struct {
	sensor.UnimplementedTelemetryServiceServer
}

func (s *grpcServer) SendData(ctx context.Context, in *sensor.TelemetryData) (*sensor.Reply, error) {
	log.Printf("[gRPC] Ricevuto da: %s, Temp: %.2f", in.SensorId, in.Temperature)
	return &sensor.Reply{Status: "OK", ServerTime: time.Now().Unix()}, nil
}

// Handler per i dati telemetry (REST)
func restHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := io.ReadAll(r.Body)
	var data sensor.TelemetryData
	json.Unmarshal(body, &data)

	log.Printf("[REST] Ricevuto da: %s, Temp: %.2f (Size: %d bytes)", data.SensorId, data.Temperature, len(body))
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // CORS per il frontend
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "OK", "size": len(body)})
}

func main() {
	// 1. Avvio gRPC su porta 50051
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		sensor.RegisterTelemetryServiceServer(s, &grpcServer{})
		log.Println("🚀 gRPC Server in ascolto su :50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// 2. Setup del Multiplexer per REST e Frontend
	mux := http.NewServeMux()

	// Endpoint per i dati dei sensori
	mux.HandleFunc("/telemetry", restHandler)

	// Endpoint per le statistiche del frontend (usato da app.js)
	mux.HandleFunc("/telemetry_stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*") 
		json.NewEncoder(w).Encode(map[string]string{
			"status": "online",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Serve i file statici (index.html e app.js) dalla cartella /web
	// Quando vai su http://localhost:8080/ vedrai la dashboard
	fs := http.FileServer(http.Dir("./web"))
	mux.Handle("/", fs)

	log.Println("🌐 REST Server e Dashboard in ascolto su :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}