package main

import "sync"

var (
    metricsMu sync.Mutex

    history          []Metric
    lastGlobalGrpcTS int64

    countSizeRest float64
    countSizeGrpc float64

    totalPayloadRest  int64
    totalOverheadRest int64
    
    totalPayloadGrpc  int64
    totalOverheadGrpc int64

    currentMode = "polling"
    currentSize = "small"
)