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
    initialHeaderSize := calculateGRPCOverhead(md)
    isFirstMessage := true

    for {
        in, err := stream.Recv()
        if err != nil {
            return err
        }

        // 1. Calcolo latenza IMMEDIATO per massima precisione
        lat := calculateLatency(in.Timestamp)

        // 2. Calcolo overhead
        var currentOverhead int64
        if isFirstMessage {
            currentOverhead = initialHeaderSize + 14
            isFirstMessage = false
        } else {
            currentOverhead = 14
        }

        // 3. Marshalling
        pSize, mTime := getProtoMetrics(in)
        
        // 4. Invio al Worker (Operazione non bloccante)
        metricsChan <- Metric{
            Protocol:     "gRPC",
            LatencyMs:    lat,
            PayloadByte:  pSize,
            OverheadByte: currentOverhead,
            MarshalTime:  mTime,
            Timestamp:    time.Now().Format("15:04:05.000"),
        }
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

            // 1. Prepariamo la history filtrata per gRPC
            var grpcHistory []*pb.MetricPoint
            metricsMu.Lock() 
            for _, m := range fullData.History {
                if m.Protocol == "gRPC" {
                    grpcHistory = append(grpcHistory, &pb.MetricPoint{
                        Protocol:  m.Protocol,
                        LatencyMs: m.LatencyMs,
                        Timestamp: m.Timestamp,
                        P99:       m.P99,
                    })
                }
            }
            metricsMu.Unlock()

            // 2. Inviamo l'oggetto completo di History
            grpcStats := &pb.GrpcStats{
                AvgLatency:  fullData.AvgGrpc,
                P99Latency:  fullData.P99Grpc,
                Timestamp:   fullData.LastGrpcTSRaw,
                PayloadSize: fullData.TotalPayloadGrpc,
                Overhead:    fullData.TotalOverheadGrpc,
                Throughput:  fullData.ThroughputGrpc,
                MarshalTime: fullData.MarshalAvgGrpc,
                History:     grpcHistory, 
            }

            if err := stream.Send(grpcStats); err != nil {
                return err
            }
        }
    }
}

func (s *telemetryServer) GetStats(ctx context.Context, in *pb.Empty) (*pb.GrpcStats, error) {
    fullData := getDashboardData()
    
    // Filtriamo la history per gRPC
    var grpcHistory []*pb.MetricPoint
    metricsMu.Lock()
    for _, m := range fullData.History {
        if m.Protocol == "gRPC" {
            grpcHistory = append(grpcHistory, &pb.MetricPoint{
                Protocol:  m.Protocol,
                LatencyMs: m.LatencyMs,
                Timestamp: m.Timestamp,
                P99:       m.P99,
            })
        }
    }
    metricsMu.Unlock()

    return &pb.GrpcStats{
        AvgLatency:  fullData.AvgGrpc,
        P99Latency:  fullData.P99Grpc,
        PayloadSize: fullData.TotalPayloadGrpc,
        Overhead:    fullData.TotalOverheadGrpc,
        Throughput:  fullData.ThroughputGrpc,
        MarshalTime: fullData.MarshalAvgGrpc,
        History:     grpcHistory,
    }, nil
}