package main

import "sync"

// Definizioni generiche utilizzate da tutto il programma
type Metric struct {
    Protocol     string  `json:"protocol"`
    LatencyMs    float64 `json:"latency_ms"`
    Timestamp    string  `json:"timestamp"`     // Per il JSON (Dashboard REST)
    RawTimestamp int64   `json:"raw_timestamp"` // Per il gRPC (Sincronizzazione)
}

type DashboardResponse struct {
    History []Metric `json:"history"`
    AvgRest float64  `json:"avg_rest"`
    AvgGrpc float64  `json:"avg_grpc"`
    P99Rest float64  `json:"p99_rest"`
    P99Grpc float64  `json:"p99_grpc"`
	LastGrpcTSRaw int64 `json:"last_grpc_ts_raw"` 
}

var (
    metricsMu sync.Mutex 
    history   []Metric   
	lastGlobalGrpcTS int64
)

// Funzione generica per processare i dati in ingresso
func processIncomingData(protocol string, timestamp int64) {
    saveMetric(protocol, timestamp)
}