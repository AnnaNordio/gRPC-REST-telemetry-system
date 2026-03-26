package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"fmt"

	"google.golang.org/grpc"
	pb "telemetry-bench/proto"
)

type server struct {
	pb.UnimplementedTelemetryServiceServer
}
var currentMode = "polling"

func (s *server) StreamData(stream pb.TelemetryService_StreamDataServer) error {
    for {
        in, err := stream.Recv()
        if err != nil {
            return err // Fine dello stream
        }
        processIncomingData("gRPC", in.LatencyGrpc)
    }
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
		
		response := getDashboardData()
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/set-mode", func(w http.ResponseWriter, r *http.Request) {
		newMode := r.URL.Query().Get("mode")
		if newMode != "" && newMode != currentMode {
			currentMode = newMode
			resetStats() 
		}
		
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/get-mode", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, currentMode)
	})

	// Servire la dashboard
	fs := http.FileServer(http.Dir("dashboard"))
	http.Handle("/", fs)

	log.Println("🚀 Dashboard e API disponibili su http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}