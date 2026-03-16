package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"google.golang.org/grpc"
	pb "telemetry-bench/proto"
)

type Metric struct {
	Protocol  string  `json:"protocol"`
	LatencyMs float64 `json:"latency_ms"`
	Temp      float32 `json:"temp"` // Aggiunto per la dashboard
	Timestamp string  `json:"timestamp"`
}

var (
	metricsMu sync.Mutex
	history   []Metric
)

type server struct {
	pb.UnimplementedTelemetryServiceServer
}

func saveMetric(protocol string, start time.Time, temp float32) {
	elapsed := float64(time.Since(start).Nanoseconds()) / 1000.0
	metricsMu.Lock()
	defer metricsMu.Unlock()

	history = append(history, Metric{
		Protocol:  protocol,
		LatencyMs: elapsed,
		Temp:      temp,
		Timestamp: time.Now().Format("15:04:05"),
	})
	if len(history) > 200 { // Aumentato un po' il buffer
		history = history[1:]
	}
}

func (s *server) SendData(ctx context.Context, in *pb.SensorData) (*pb.Empty, error) {
	saveMetric("gRPC", time.Now())
	return &pb.Empty{}, nil
}

func main() {
	// 1. gRPC Server
	go func() {
		lis, _ := net.Listen("tcp", ":50051")
		s := grpc.NewServer()
		pb.RegisterTelemetryServiceServer(s, &server{})
		log.Println("🚀 gRPC Server in ascolto su :50051")
		s.Serve(lis)
	}()

	// 2. Endpoint API per invio dati e risultati
	http.HandleFunc("/telemetry", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var data pb.SensorData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			return
		}
		saveMetric("REST", start, data.Temperature)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		metricsMu.Lock()
		json.NewEncoder(w).Encode(history)
		metricsMu.Unlock()
	})

	fs := http.FileServer(http.Dir("dashboard"))
	http.Handle("/", fs) 

	log.Println("🚀 Dashboard e API disponibili su http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}