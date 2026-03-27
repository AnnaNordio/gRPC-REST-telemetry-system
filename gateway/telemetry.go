package main

import (
    "fmt"
    "time"
)

func saveMetric(protocol string, sensorTimestamp int64) {
    if sensorTimestamp <= 0 { return }

    metricsMu.Lock()
    defer metricsMu.Unlock()

    if protocol == "gRPC" {
        if grpcCount < warmupThreshold {
            grpcCount++
            fmt.Printf("⏳ [Warm-up gRPC] %d/%d\n", grpcCount, warmupThreshold)
            return
        }
        lastGlobalGrpcTS = sensorTimestamp
    } else if protocol == "REST" {
        if restCount < warmupThreshold {
            restCount++
            fmt.Printf("⏳ [Warm-up REST] %d/%d\n", restCount, warmupThreshold)
            return
        }
    }
    // ---------------------------

    now := time.Now().UnixMicro()
    realLatency := float64(now-sensorTimestamp) 
    displayTS := time.UnixMicro(sensorTimestamp).Format("15:04:05.000") 

    history = append(history, Metric{
        Protocol:     protocol,
        LatencyMs:    realLatency, 
        Timestamp:    displayTS,
        RawTimestamp: sensorTimestamp,
    })

    if len(history) > 200 {
        history = history[1:]
    }
}

func savePayload(protocol string, size int) {
    metricsMu.Lock()
    defer metricsMu.Unlock()

    if protocol == "REST" {
        if restCount < warmupThreshold { return }
        sumSizeRest += float64(size)
        countSizeRest++
    } else {
        if grpcCount < warmupThreshold { return }
        sumSizeGrpc += float64(size)
        countSizeGrpc++
    }
}

func resetStats() {
    metricsMu.Lock()
    defer metricsMu.Unlock()
    
    // Azzeriamo tutto: history, medie e contatori di warm-up
    history = []Metric{}
    sumSizeRest = 0
    countSizeRest = 0
    sumSizeGrpc = 0
    countSizeGrpc = 0
    
    fmt.Println("Reset completo: statistiche azzerate.")
}