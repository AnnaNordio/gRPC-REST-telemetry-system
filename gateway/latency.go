package main

import (
	"time"
	"sort"
)

type Metric struct {
	Protocol  string  `json:"protocol"`
	LatencyMs float64 `json:"latency_ms"`
	Timestamp string  `json:"timestamp"`
}

var (
	history            []Metric
	sumRest, countRest float64
	sumGrpc, countGrpc float64
)

func saveMetric(protocol string, latency float64) {
	metricsMu.Lock()
	defer metricsMu.Unlock()

	if protocol == "REST" {
		sumRest += latency
		countRest++
	} else {
		sumGrpc += latency
		countGrpc++
	}

	history = append(history, Metric{
		Protocol:  protocol,
		LatencyMs: latency,
		Timestamp: time.Now().Format("15:04:05"),
	})

	if len(history) > 200 {
		history = history[1:]
	}
}

func getDashboardData() DashboardResponse {
    metricsMu.Lock()
    defer metricsMu.Unlock()

    var restLats, grpcLats []float64
    for _, m := range history {
        if m.Protocol == "REST" {
            restLats = append(restLats, m.LatencyMs)
        } else {
            grpcLats = append(grpcLats, m.LatencyMs)
        }
    }

    avgR, avgG := 0.0, 0.0
    if countRest > 0 { avgR = sumRest / countRest }
    if countGrpc > 0 { avgG = sumGrpc / countGrpc }

    return DashboardResponse{
        History: history,
        AvgRest: avgR,
        AvgGrpc: avgG,
        P99Rest: calculatePercentile(restLats, 0.99),
        P99Grpc: calculatePercentile(grpcLats, 0.99),
    }
}

func resetStats() {
    metricsMu.Lock()
    defer metricsMu.Unlock()

    history = []Metric{}

    sumRest = 0
    countRest = 0
    sumGrpc = 0
    countGrpc = 0

}

func calculatePercentile(latencies []float64, percentile float64) float64 {
    if len(latencies) == 0 {
        return 0
    }
    // Creiamo una copia per non sporcare i dati originali
    sorted := make([]float64, len(latencies))
    copy(sorted, latencies)
    sort.Float64s(sorted)

    // Calcoliamo l'indice (N * percentile)
    index := int(float64(len(sorted)-1) * percentile)
    return sorted[index]
}