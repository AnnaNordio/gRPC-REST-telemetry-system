package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	pb "telemetry-bench/proto"
)

// --- CONFIGURAZIONE E STATO ---

type server struct {
	pb.UnimplementedTelemetryServiceServer
}

var currentMode = "polling"
var currentSize = "small"

// --- MIDDLEWARE CORS ---
// Questa funzione avvolge i tuoi handler e aggiunge i permessi per il browser
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permetti l'origine del Frontend (o tutte con *)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Gestione richiesta "Preflight" (il browser la invia prima del POST)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// --- IMPLEMENTAZIONE gRPC ---

func (s *server) StreamData(stream pb.TelemetryService_StreamDataServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			return err
		}
		saveMetric("gRPC", in.Timestamp)
	}
}

func (s *server) SendData(ctx context.Context, in *pb.SensorData) (*pb.Empty, error) {
	saveMetric("gRPC", in.Timestamp)
	return &pb.Empty{}, nil
}

func (s *server) GetGrpcStream(in *pb.Empty, stream pb.TelemetryService_GetGrpcStreamServer) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-ticker.C:
			fullData := getDashboardData()
			grpcStats := &pb.GrpcStats{
				AvgLatency: fullData.AvgGrpc,
				P99Latency: fullData.P99Grpc,
				Timestamp:  fullData.LastGrpcTSRaw,
			}
			if err := stream.Send(grpcStats); err != nil {
				return err
			}
		}
	}
}

func (s *server) GetStats(ctx context.Context, in *pb.Empty) (*pb.GrpcStats, error) {
	fullData := getDashboardData()
	return &pb.GrpcStats{
		AvgLatency: fullData.AvgGrpc,
		P99Latency: fullData.P99Grpc,
		Timestamp:  fullData.LastGrpcTSRaw,
	}, nil
}

// --- MAIN E HTTP SERVER ---

func main() {
	// 1. Avvio gRPC Server in background
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterTelemetryServiceServer(s, &server{})
		log.Println("🚀 gRPC Server in ascolto su :50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// 2. Definizione del Multiplexer (Mux) per HTTP
	mux := http.NewServeMux()

	// Handler per i dati della dashboard
	mux.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fullData := getDashboardData()

		var restOnlyHistory []Metric
		for _, m := range fullData.History {
			if m.Protocol == "REST" {
				restOnlyHistory = append(restOnlyHistory, m)
			}
		}

		response := struct {
			History []Metric `json:"history"`
			AvgRest float64  `json:"avg_rest"`
			P99Rest float64  `json:"p99_rest"`
		}{
			History: restOnlyHistory,
			AvgRest: fullData.AvgRest,
			P99Rest: fullData.P99Rest,
		}
		json.NewEncoder(w).Encode(response)
	})

	// Handler per ricevere telemetria REST dai sensori
	mux.HandleFunc("/telemetry", func(w http.ResponseWriter, r *http.Request) {
		var data pb.SensorData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		saveMetric("REST", data.Timestamp)
		w.WriteHeader(http.StatusOK)
	})

	// Configurazione Modalità
	mux.HandleFunc("/set-mode", func(w http.ResponseWriter, r *http.Request) {
		newMode := r.URL.Query().Get("mode")
		if newMode != "" && newMode != currentMode {
			currentMode = newMode
			resetStats()
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/get-mode", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, currentMode)
	})

	// Configurazione Dimensione Payload
	mux.HandleFunc("/set-size", func(w http.ResponseWriter, r *http.Request) {
		newSize := r.URL.Query().Get("size")
		if newSize != "" && newSize != currentSize {
			currentSize = newSize
			resetStats()
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/get-size", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, currentSize)
	})

	// File statici della Dashboard
	fs := http.FileServer(http.Dir("dashboard"))
	mux.Handle("/", fs)

	// 3. Avvio del server HTTP con Middleware CORS applicato
	log.Println("🚀 Gateway e Dashboard attivi su http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(mux)))
}