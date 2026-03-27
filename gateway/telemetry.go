package main

import (
    "time"
)

/*func saveMetric(protocol string, sensorTimestamp int64) {
    if sensorTimestamp <= 0 { return }

    metricsMu.Lock()
    defer metricsMu.Unlock()

    if protocol == "gRPC" {
        if grpcCount < warmupThreshold {
            grpcCount++
            fmt.Printf("⏳ [Warm-up gRPC] %d/%d\n", grpcCount, warmupThreshold)
            return
        }
        lastGlobalGrpcTS = sensorTimestamp
    } else if protocol == "REST" {
        if restCount < warmupThreshold {
            restCount++
            fmt.Printf("⏳ [Warm-up REST] %d/%d\n", restCount, warmupThreshold)
            return
        }
    }
    // ---------------------------

    now := time.Now().UnixMicro()
    realLatency := float64(now-sensorTimestamp) 
    displayTS := time.UnixMicro(sensorTimestamp).Format("15:04:05.000") 

    history = append(history, Metric{
        Protocol:     protocol,
        LatencyMs:    realLatency, 
        Timestamp:    displayTS,
        RawTimestamp: sensorTimestamp,
    })

    if len(history) > 200 {
        history = history[1:]
    }
}

func savePayload(protocol string, pSize int64, hSize int64) {
    metricsMu.Lock()
    defer metricsMu.Unlock()

    // Incremento contatori per il warmup
    if protocol == "REST" {
        restCount++
        if restCount <= warmupThreshold { return }
    } else {
        grpcCount++
        if grpcCount <= warmupThreshold { return }
    }

    // Aggiungiamo il record alla history per i grafici temporali
    newMetric := Metric{
        Protocol:     protocol,
        PayloadByte:  pSize,
        OverheadByte: hSize,
        Timestamp:    string(time.Now().UnixMilli()),
    }
    
    // Manteniamo la history pulita (es. ultimi 100 elementi)
    if len(history) > 100 {
        history = history[1:]
    }
    history = append(history, newMetric)

    // Aggiorniamo comunque le somme globali per le medie generali
    if protocol == "REST" {
        sumSizeRest += float64(pSize)
        sumOverheadRest += float64(hSize)
        countSizeRest++
    } else {
        sumSizeGrpc += float64(pSize)
        sumOverheadGrpc += float64(hSize)
        countSizeGrpc++
    }
}*/

func saveAllMetrics(protocol string, sensorTS int64, pSize int64, hSize int64) {
    if sensorTS <= 0 { return }

    metricsMu.Lock()
    defer metricsMu.Unlock()

    // 1. Warmup Check
    if protocol == "gRPC" {
        if grpcCount < warmupThreshold {
            grpcCount++; return
        }
        lastGlobalGrpcTS = sensorTS
    } else {
        if restCount < warmupThreshold {
            restCount++; return
        }
    }

    // 2. Calcolo Latenza
    now := time.Now().UnixMicro()
    latency := float64(now - sensorTS)

    // 3. Aggiornamento Totali Cumulativi (quelli che usi per le card)
    if protocol == "REST" {
        sumSizeRest += float64(pSize)
        sumOverheadRest += float64(hSize)
    } else {
        sumSizeGrpc += float64(pSize)
        sumOverheadGrpc += float64(hSize)
    }

    // 4. Aggiunta alla History (UNICA ENTRY)
    history = append(history, Metric{
        Protocol:     protocol,
        LatencyMs:    latency,
        PayloadByte:  pSize,
        OverheadByte: hSize,
        Timestamp:    time.UnixMicro(sensorTS).Format("15:04:05.000"),
        RawTimestamp: sensorTS,
    })

    if len(history) > 200 {
        history = history[1:]
    }
}

func resetStats() {
    metricsMu.Lock()
    defer metricsMu.Unlock()
    
    // Azzeriamo tutto: history, medie e contatori di warm-up
    history = []Metric{}
    sumSizeRest = 0
    countSizeRest = 0
    sumSizeGrpc = 0
    countSizeGrpc = 0
	grpcCount = 0
    restCount = 0
    
}