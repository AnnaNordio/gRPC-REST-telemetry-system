package main

type Metric struct {
    Protocol    string  `json:"protocol"`
    LatencyMs   float64 `json:"latency_ms"`
    PayloadByte int64   `json:"payload_byte"`  
    OverheadByte int64  `json:"overhead_byte"` 
    Timestamp   string   `json:"timestamp"`
	RawTimestamp int64   `json:"raw_timestamp"`
}

type DashboardResponse struct {
    History         []Metric `json:"history"`
    AvgRest         float64  `json:"avg_rest"`
    AvgGrpc         float64  `json:"avg_grpc"`
    P99Rest         float64  `json:"p99_rest"`
    P99Grpc         float64  `json:"p99_grpc"`
    TotalRestSize     float64 `json:"total_rest_size"`
    TotalRestOverhead float64 `json:"total_rest_overhead"`
    TotalGrpcSize     float64 `json:"total_grpc_size"`
    TotalGrpcOverhead float64 `json:"total_grpc_overhead"`
    LastGrpcTSRaw   int64    `json:"last_grpc_ts_raw"`
}