package main

import (
	"encoding/csv"
	"os"
)

type Metric struct {
	Protocol     string  `json:"protocol"`
	LatencyMs    float64 `json:"latency_ms"`
	Timestamp    string  `json:"timestamp"`
	PayloadByte  int64   `json:"payload_byte"`
	OverheadByte int64   `json:"overhead_byte"`
	P99          float64 `json:"p99_ms"`
	MarshalTime  float64 `json:"marshal_time_us"`
	Throughput   float64 `json:"throughput_bytes_per_sec"`
}

type DashboardResponse struct {
	History           []Metric `json:"history"`
	AvgRest           float64  `json:"avg_rest"`
	AvgGrpc           float64  `json:"avg_grpc"`
	P99Rest           float64  `json:"p99_rest"`
	P99Grpc           float64  `json:"p99_grpc"`
	MarshalAvgRest    float64  `json:"marshal_avg_rest"`
	MarshalAvgGrpc    float64  `json:"marshal_avg_grpc"`
	ThroughputGrpc    float64  `json:"throughput_grpc"`
	ThroughputRest    float64  `json:"throughput_rest"`
	TotalPayloadRest  int64    `json:"total_payload_rest"`
	TotalOverheadRest int64    `json:"total_overhead_rest"`
	TotalPayloadGrpc  int64    `json:"total_payload_grpc"`
	TotalOverheadGrpc int64    `json:"total_overhead_grpc"`
	LastGrpcTSRaw     int64    `json:"last_grpc_ts_raw"`
}

type MetricsWriter struct {
	file       *os.File
	csvWriter  *csv.Writer
	lastConfig string
}
