package main

import (
    "time"
)

func saveMetric(protocol string, sensorTimestamp int64) {
    if sensorTimestamp <= 0 { return }

    metricsMu.Lock()
    defer metricsMu.Unlock()

    if protocol == "gRPC" {
        lastGlobalGrpcTS = sensorTimestamp
    }

    now := time.Now().UnixMicro()
    realLatency := float64(now - sensorTimestamp)
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
        sumSizeRest += float64(size)
        countSizeRest++
    } else {
        sumSizeGrpc += float64(size)
        countSizeGrpc++
    }
}

func resetStats() {
    metricsMu.Lock()
    defer metricsMu.Unlock()
    history = []Metric{}
}