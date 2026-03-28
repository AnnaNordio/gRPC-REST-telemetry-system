package main

import (
    "context"
    "time"
    pb "telemetry-bench/proto"
)

type telemetryServer struct {
    pb.UnimplementedTelemetryServiceServer
}

func (s *telemetryServer) StreamData(stream pb.TelemetryService_StreamDataServer) error {
    for {
        in, err := stream.Recv()
        if err != nil {
            return err
        }
        saveMetric("gRPC", in.Timestamp)
    }
}

func (s *telemetryServer) SendData(ctx context.Context, in *pb.SensorData) (*pb.Empty, error) {
    saveMetric("gRPC", in.Timestamp)
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
            }
            if err := stream.Send(grpcStats); err != nil {
                return err
            }
        }
    }
}

func (s *telemetryServer) GetStats(ctx context.Context, in *pb.Empty) (*pb.GrpcStats, error) {
    fullData := getDashboardData()
    return &pb.GrpcStats{
        AvgLatency: fullData.AvgGrpc,
        P99Latency: fullData.P99Grpc,
        Timestamp:  fullData.LastGrpcTSRaw,
    }, nil
}