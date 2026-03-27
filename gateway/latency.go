package main

import (
    "sort"
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

// Calcola tutti i dati aggregati per la Dashboard
func getDashboardData() DashboardResponse {
    metricsMu.Lock()
    defer metricsMu.Unlock()

    var restLats, grpcLats []float64
    var sumR, sumG float64

    for _, m := range history {
        if m.Protocol == "REST" {
            restLats = append(restLats, m.LatencyMs)
            sumR += m.LatencyMs
        } else {
            grpcLats = append(grpcLats, m.LatencyMs)
            sumG += m.LatencyMs
        }
    }

    return DashboardResponse{
        History:       history,
        AvgRest:       safeAvg(sumR, len(restLats)),
        AvgGrpc:       safeAvg(sumG, len(grpcLats)),
        P99Rest:       calculatePercentile(restLats, 0.99),
        P99Grpc:       calculatePercentile(grpcLats, 0.99),
        LastGrpcTSRaw: lastGlobalGrpcTS,
    }
}

// Calcolo matematico della media
func safeAvg(sum float64, count int) float64 {
    if count == 0 {
        return 0
    }
    return sum / float64(count)
}

// Calcolo matematico del percentile (P99)
func calculatePercentile(latencies []float64, percentile float64) float64 {
    if len(latencies) == 0 {
        return 0
    }
    
    // Copia e ordina per non alterare la history originale
    sorted := make([]float64, len(latencies))
    copy(sorted, latencies)
    sort.Float64s(sorted)

    index := int(float64(len(sorted)-1) * percentile)
    return sorted[index]
}

// Pulisce la history
func resetStats() {
    metricsMu.Lock()
    defer metricsMu.Unlock()
    history = []Metric{}
}