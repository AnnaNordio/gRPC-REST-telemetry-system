package main

import (
    "sort"
    "time"
)


func saveMetric(protocol string, latency float64, sensorTimestamp string) {
    metricsMu.Lock()
    defer metricsMu.Unlock()

    // Se il sensore non ha mandato il timestamp (es. polling vecchio), usa quello attuale
    ts := sensorTimestamp
    if ts == "" {
        ts = time.Now().Format("15:04:05")
    }

    history = append(history, Metric{
        Protocol:  protocol,
        LatencyMs: latency,
        Timestamp: ts,
    })

    if len(history) > 200 {
        history = history[1:len(history)]
    }
}

// Calcola tutti i dati aggregati per la Dashboard
func getDashboardData() DashboardResponse {
    metricsMu.Lock() // RLock non disponibile con sync.Mutex semplice, usiamo Lock
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
        History: history,
        AvgRest: safeAvg(sumR, len(restLats)),
        AvgGrpc: safeAvg(sumG, len(grpcLats)),
        P99Rest: calculatePercentile(restLats, 0.99),
        P99Grpc: calculatePercentile(grpcLats, 0.99),
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