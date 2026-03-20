package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	pb "telemetry-bench/proto"
)

type server struct {
	pb.UnimplementedTelemetryServiceServer
}

func (s *server) SendData(ctx context.Context, in *pb.SensorData) (*pb.Empty, error) {
	processIncomingData("gRPC", in.LatencyGrpc)
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

	// 2. HTTP Handlers
	http.HandleFunc("/telemetry", func(w http.ResponseWriter, r *http.Request) {
		var data pb.SensorData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			return
		}
		processIncomingData("REST", data.LatencyRest)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		// Usiamo la funzione definita in latency.go
		response := getDashboardData()
		json.NewEncoder(w).Encode(response)
	})

	// Servire la dashboard
	fs := http.FileServer(http.Dir("dashboard"))
	http.Handle("/", fs)

	log.Println("🚀 Dashboard e API disponibili su http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}