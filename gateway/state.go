package main

import (
    "sync"
    "time"   
    "telemetry-bench/pkg/config"
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

    activeConfig = config.TelemetryConfig{
        Mode:     "polling",
        Size:     "small",
        Sensors:  1,
        Protocol: "both",
    }

    msgCountRest uint64
    msgCountGrpc uint64
    
    throughputRest float64
    throughputGrpc float64

    warmupUntil time.Time
    warmupDuration = 30 * time.Second
)