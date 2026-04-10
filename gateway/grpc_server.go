package main

import (
    "context"
    "time"
    "google.golang.org/grpc/metadata"  
    pb "telemetry-bench/proto"
)

type telemetryServer struct {
    pb.UnimplementedTelemetryServiceServer
}

func (s *telemetryServer) StreamData(stream pb.TelemetryService_StreamDataServer) error {
    md, _ := metadata.FromIncomingContext(stream.Context())
    
    // Calcoliamo l'overhead iniziale (Headers/Metadata)
    initialHeaderSize := calculateGRPCOverhead(md)
    
    isFirstMessage := true

    for {
        in, err := stream.Recv()
        if err != nil {
            return err
        }

        var currentOverhead int64
        if isFirstMessage {
            // Primo messaggio: Headers + Frame gRPC (14)
            currentOverhead = initialHeaderSize + 14
            isFirstMessage = false
        } else {
            // Messaggi successivi: SOLO il frame gRPC (14)
            currentOverhead = 14
        }

        pSize := int64(getProtoSize(in))
        
        // Invia il peso del SINGOLO evento
        processAndStoreMetrics("gRPC", in, pSize, currentOverhead)
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
                Throughput: fullData.ThroughputGrpc,
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
        PayloadSize: fullData.TotalPayloadGrpc,
        Overhead: fullData.TotalOverheadGrpc,
        Throughput: fullData.ThroughputGrpc,
    }, nil
}