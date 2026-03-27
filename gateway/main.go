package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"fmt"
	"time"

	"google.golang.org/grpc"
	pb "telemetry-bench/proto"
)

type server struct {
	pb.UnimplementedTelemetryServiceServer
}
var currentMode = "polling"
var currentSize = "small"


func (s *server) StreamData(stream pb.TelemetryService_StreamDataServer) error {
    for {
        in, err := stream.Recv()
        if err != nil { return err }
        // Salva solo il dato gRPC
        saveMetric("gRPC", in.LatencyGrpc, string(in.Timestamp))
    }
}

func (s *server) SendData(ctx context.Context, in *pb.SensorData) (*pb.Empty, error) {
	processIncomingData("gRPC", in.LatencyGrpc, string(in.Timestamp))
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
    }, nil
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
		json.NewDecoder(r.Body).Decode(&data)
		// Salva solo il dato REST
		saveMetric("REST", data.LatencyRest, string(data.Timestamp))
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

	http.HandleFunc("/set-size", func(w http.ResponseWriter, r *http.Request) {
		newSize := r.URL.Query().Get("size")
		if newSize != "" && newSize != currentSize {
			currentSize = newSize
		}
		
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/get-size", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, currentSize)
	})

	// Servire la dashboard
	fs := http.FileServer(http.Dir("dashboard"))
	http.Handle("/", fs)

	log.Println("🚀 Dashboard e API disponibili su http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}