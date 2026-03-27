package main

import (
	"context"
	"fmt"
	"time"

	pb "telemetry-bench/proto"
)

type telemetryServer struct {
	pb.UnimplementedTelemetryServiceServer
}

// 1. STREAMING: Ricezione dati dal sensore
// L'overhead e il payload vengono gestiti dallo STREAM INTERCEPTOR
func (s *telemetryServer) StreamData(stream pb.TelemetryService_StreamDataServer) error {
	for {
		_, err := stream.Recv() // Riceviamo e basta
		if err != nil {
			return err
		}
		// NON CHIAMARE NULLA QUI. L'interceptor fa tutto.
	}
}

// 2. UNARY: Ricezione dati singoli dal sensore
// L'overhead e il payload vengono gestiti dallo UNARY INTERCEPTOR
func (s *telemetryServer) SendData(ctx context.Context, in *pb.SensorData) (*pb.Empty, error) {
	// NON CHIAMARE NULLA QUI. L'interceptor fa tutto.
	return &pb.Empty{}, nil
}

// 3. STREAMING VERSO DASHBOARD: Invio statistiche aggiornate
// Qui DOBBIAMO convertire i tipi int64 in float64 per il proto
func (s *telemetryServer) GetGrpcStream(in *pb.Empty, stream pb.TelemetryService_GetGrpcStreamServer) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case <-ticker.C:
			fullData := getDashboardData()
			grpcStats := &pb.GrpcStats{
				AvgLatency:  fullData.AvgGrpc,
				P99Latency:  fullData.P99Grpc,
                Timestamp:   fullData.LastGrpcTSRaw,
				PayloadSize: float64(fullData.TotalGrpcSize),
				Overhead:    float64(fullData.TotalGrpcOverhead),
			}
			if err := stream.Send(grpcStats); err != nil {
				return err
			}
		}
	}
}

// 4. UNARY VERSO DASHBOARD: Richiesta statistiche singola
func (s *telemetryServer) GetStats(ctx context.Context, in *pb.Empty) (*pb.GrpcStats, error) {
	fullData := getDashboardData()
	
	fmt.Printf("📊 [Stats Request] Lat:%.2f Payload:%d\n", fullData.TotalGrpcOverhead, fullData.TotalGrpcSize)
	
	return &pb.GrpcStats{
		AvgLatency:  fullData.AvgGrpc,
		P99Latency:  fullData.P99Grpc,
        Timestamp:   fullData.LastGrpcTSRaw,
		PayloadSize: float64(fullData.TotalGrpcSize),
		Overhead:    float64(fullData.TotalGrpcOverhead),
	}, nil
}