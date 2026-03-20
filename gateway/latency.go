package main

import (
	"time"
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

	avgR, avgG := 0.0, 0.0
	if countRest > 0 { avgR = sumRest / countRest }
	if countGrpc > 0 { avgG = sumGrpc / countGrpc }

	// I nomi dei campi qui sotto devono essere IDENTICI a metrics.go
	return DashboardResponse{
		History: history,
		AvgRest: avgR,
		AvgGrpc: avgG,
	}
}