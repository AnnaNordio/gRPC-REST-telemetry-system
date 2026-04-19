package main

import (
    "sync"
    "time"   
)


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
    currentSensors = "1"
    currentProtocol = "both"

    msgCountRest uint64
    msgCountGrpc uint64
    
    throughputRest float64
    throughputGrpc float64

    warmupUntil time.Time
    warmupDuration = 5 * time.Second
)