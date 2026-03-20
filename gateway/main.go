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
	Timestamp string  `json:"timestamp"`
}

var (
	metricsMu sync.Mutex
	history   []Metric
	sumRest, countRest float64
    sumGrpc, countGrpc float64
)

type DashboardResponse struct {
    History []Metric `json:"history"`
    AvgRest float64  `json:"avg_rest"`
    AvgGrpc float64  `json:"avg_grpc"`
}

type server struct {
	pb.UnimplementedTelemetryServiceServer
}

func saveMetric(protocol string, latency float64, temp float32) {
    metricsMu.Lock()
    defer metricsMu.Unlock()

    if protocol == "REST" {
        sumRest += latency
        countRest++
    } else {
        sumGrpc += latency
        countGrpc++
    }

    history = append(history, Metric{
        Protocol:  protocol,
        LatencyMs: latency,
        Timestamp: time.Now().Format("15:04:05"),
    })
    
    if len(history) > 200 {
        history = history[1:]
    }
}

func (s *server) SendData(ctx context.Context, in *pb.SensorData) (*pb.Empty, error) {
	saveMetric("gRPC", in.LatencyGrpc, in.Temperature)
	return &pb.Empty{}, nil
}

func main() {
	go func() {
		lis, _ := net.Listen("tcp", ":50051")
		s := grpc.NewServer()
		pb.RegisterTelemetryServiceServer(s, &server{})
		log.Println("🚀 gRPC Server in ascolto su :50051")
		s.Serve(lis)
	}()

	// 2. Endpoint API per invio dati e risultati
	http.HandleFunc("/telemetry", func(w http.ResponseWriter, r *http.Request) {
		var data pb.SensorData
		json.NewDecoder(r.Body).Decode(&data)
		saveMetric("REST", data.LatencyRest, data.Temperature)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		metricsMu.Lock()
		avgR, avgG := 0.0, 0.0
		if countRest > 0 { avgR = sumRest / countRest }
		if countGrpc > 0 { avgG = sumGrpc / countGrpc }
		
		response := DashboardResponse{
			History: history,
			AvgRest: avgR,
			AvgGrpc: avgG,
		}
		json.NewEncoder(w).Encode(response)
		metricsMu.Unlock()
	})

	fs := http.FileServer(http.Dir("dashboard"))
	http.Handle("/", fs) 

	log.Println("🚀 Dashboard e API disponibili su http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}