package main

import "sync"

var (
	metricsMu sync.Mutex
)

// Definiamo qui la struct globale
type DashboardResponse struct {
	History []Metric `json:"history"`
	AvgRest float64  `json:"avg_rest"`
	AvgGrpc float64  `json:"avg_grpc"`
}

// Funzione di smistamento
func processIncomingData(protocol string, latency float64) {
	saveMetric(protocol, latency)
}