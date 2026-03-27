package main

import "sync"

// Definizioni generiche utilizzate da tutto il programma
type Metric struct {
    Protocol  string  `json:"protocol"`
    LatencyMs float64 `json:"latency_ms"`
    Timestamp string  `json:"timestamp"`
}

type DashboardResponse struct {
    History []Metric `json:"history"`
    AvgRest float64  `json:"avg_rest"`
    AvgGrpc float64  `json:"avg_grpc"`
    P99Rest float64  `json:"p99_rest"`
    P99Grpc float64  `json:"p99_grpc"`
}

var (
    metricsMu sync.Mutex // Unico Mutex globale per la coerenza dei dati
    history   []Metric   // Lo storage dei campioni
)

// Funzione generica per processare i dati in ingresso
func processIncomingData(protocol string, latency float64, timestamp string) {
    saveMetric(protocol, latency, timestamp)
}