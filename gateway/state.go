package main

import "sync"

var (
    metricsMu sync.Mutex

    history          []Metric
    lastGlobalGrpcTS int64

    sumSizeRest, countSizeRest float64
    sumSizeGrpc, countSizeGrpc float64

    currentMode = "polling"
    currentSize = "small"

	warmupThreshold = 5
    grpcCount       = 0
    restCount       = 0
)