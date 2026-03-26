package main

import "sync"

var (
	metricsMu sync.Mutex
)

type DashboardResponse struct {
	History []Metric `json:"history"`
	AvgRest float64  `json:"avg_rest"`
	AvgGrpc float64  `json:"avg_grpc"`
	P99Rest float64  `json:"p99_rest"` 
    P99Grpc float64  `json:"p99_grpc"`
}

func processIncomingData(protocol string, latency float64) {
	saveMetric(protocol, latency)
}