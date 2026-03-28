package main

import (
    "context"
    "time"
    "fmt"
    "google.golang.org/grpc/metadata"  
    pb "telemetry-bench/proto"
)

type telemetryServer struct {
    pb.UnimplementedTelemetryServiceServer
}

func (s *telemetryServer) StreamData(stream pb.TelemetryService_StreamDataServer) error {
    md, _ := metadata.FromIncomingContext(stream.Context())
    for {
        in, err := stream.Recv()
        if err != nil {
            return err
        }
        SaveGrpcMetrics(in, md)
    }
}

func (s *telemetryServer) SendData(ctx context.Context, in *pb.SensorData) (*pb.Empty, error) {
    md, _ := metadata.FromIncomingContext(ctx)
    SaveGrpcMetrics(in, md)
    return &pb.Empty{}, nil
}

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
                AvgLatency: fullData.AvgGrpc,
                P99Latency: fullData.P99Grpc,
                Timestamp:  fullData.LastGrpcTSRaw,
                PayloadSize: fullData.TotalPayloadGrpc,
                Overhead: fullData.TotalOverheadGrpc,
            }
            if err := stream.Send(grpcStats); err != nil {
                return err
            }
        }
    }
}

func (s *telemetryServer) GetStats(ctx context.Context, in *pb.Empty) (*pb.GrpcStats, error) {
    fullData := getDashboardData()
    fmt.Printf("gRPC Stats Requested: Avg=%.2f ms, P99=%.2f ms, Payload=%d B, Overhead=%d B\n", fullData.AvgGrpc, fullData.P99Grpc, fullData.TotalPayloadGrpc, fullData.TotalOverheadGrpc)
    return &pb.GrpcStats{
        AvgLatency: fullData.AvgGrpc,
        P99Latency: fullData.P99Grpc,
        Timestamp:  fullData.LastGrpcTSRaw,
        PayloadSize: fullData.TotalPayloadGrpc,
        Overhead: fullData.TotalOverheadGrpc,
    }, nil
}