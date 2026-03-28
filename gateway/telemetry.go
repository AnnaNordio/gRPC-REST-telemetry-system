package main

import (
    "time"
    "net/http"
    "google.golang.org/grpc/metadata"  
    pb "telemetry-bench/proto"
)

func SaveRestMetrics(data *pb.SensorData, r *http.Request) {
    pSize := int64(getJsonSize(data))
    hSize := calculateHTTPOverhead(r)
    
    processAndStoreMetrics("REST", data, pSize, hSize)
}

func SaveGrpcMetrics(data *pb.SensorData, md metadata.MD) {
    pSize := int64(getProtoSize(data))
    hSize := 5 + calculateGRPCOverhead(md) 
    
    processAndStoreMetrics("gRPC", data, pSize, hSize)
}

func processAndStoreMetrics(protocol string, data *pb.SensorData, pSize, hSize int64) {
    if data.Timestamp <= 0 {
        return
    }

    latency := calculateLatency(data.Timestamp)

    metricsMu.Lock()
    defer metricsMu.Unlock()

    // 1. Logica di Warmup
    if protocol == "gRPC" {
        if grpcCount < warmupThreshold {
            grpcCount++
            return
        }
        lastGlobalGrpcTS = data.Timestamp
        totalPayloadGrpc += pSize
        totalOverheadGrpc += hSize
    } else {
        if restCount < warmupThreshold {
            restCount++
            return
        }
        totalPayloadRest += pSize
        totalOverheadRest += hSize
    }

    // 2. Archiviazione nella History
    newMetric := Metric{
        Protocol:     protocol,
        LatencyMs:    latency,
        PayloadByte:  pSize,
        OverheadByte: hSize,
        Timestamp:    string(time.UnixMicro(data.Timestamp).Format("15:04:05.000")),
        RawTimestamp: data.Timestamp,
    }

    history = append(history, newMetric)
    if len(history) > 200 {
        history = history[1:]
    }
}

func resetStats() {
    metricsMu.Lock()
    defer metricsMu.Unlock()
    
    history = []Metric{}
    countSizeRest = 0
    countSizeGrpc = 0
    totalPayloadRest = 0
    totalOverheadRest = 0
    totalPayloadGrpc = 0
    totalOverheadGrpc = 0
	grpcCount = 0
    restCount = 0
}