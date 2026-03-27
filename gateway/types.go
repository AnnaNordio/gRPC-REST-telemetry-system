package main

type Metric struct {
    Protocol     string  `json:"protocol"`
    LatencyMs    float64 `json:"latency_ms"`
    Timestamp    string  `json:"timestamp"`     
    RawTimestamp int64   `json:"raw_timestamp"`
}

type DashboardResponse struct {
    History       []Metric `json:"history"`
    AvgRest       float64  `json:"avg_rest"`
    AvgGrpc       float64  `json:"avg_grpc"`
    P99Rest       float64  `json:"p99_rest"`
    P99Grpc       float64  `json:"p99_grpc"`
    AvgRestSize   float64  `json:"avg_rest_size"`
    AvgGrpcSize   float64  `json:"avg_grpc_size"`
    LastGrpcTSRaw int64    `json:"last_grpc_ts_raw"`
}