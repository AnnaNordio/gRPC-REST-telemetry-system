package main

import (
    "time"
    "net/http"
    "google.golang.org/grpc/metadata"  
    "sync/atomic"
    pb "telemetry-bench/proto"
)

// Canale per trasportare le metriche dai server al worker
// Capacità 10.000 per gestire picchi di 100 sensori a 10Hz
var metricsChan = make(chan Metric, 10000)

// Inizializza il worker all'avvio
func init() {
    go metricsWorker()
}

// Il Worker: l'unico punto che scrive nella history, eliminando la contesa del Mutex
func metricsWorker() {
    for m := range metricsChan {
        if time.Now().Before(warmupUntil) {
            continue 
        }
        metricsMu.Lock()
        
        if m.Protocol == "gRPC" {
            atomic.AddUint64(&msgCountGrpc, 1)
            // Usiamo il valore grezzo del timestamp se necessario per throughput
            totalPayloadGrpc += m.PayloadByte
            totalOverheadGrpc += m.OverheadByte
        } else {
            atomic.AddUint64(&msgCountRest, 1)
            totalPayloadRest += m.PayloadByte
            totalOverheadRest += m.OverheadByte
        }

        history = append(history, m)
        if len(history) > 1000 {
            history = history[1:]
        }
        
        metricsMu.Unlock()
    }
}

func SaveRestMetrics(data *pb.SensorData, r *http.Request) {

    pSize, mTime := getJsonMetrics(data)
    hSize := calculateHTTPOverhead(r)
    lat := calculateLatency(data.Timestamp)
    
    // Invia al worker invece di bloccare qui
    metricsChan <- Metric{
        Protocol:     "REST",
        LatencyMs:    lat,
        PayloadByte:  pSize,
        OverheadByte: hSize,
        MarshalTime:  mTime,
        Timestamp:    time.Now().Format("15:04:05.000"),
    }
}

func SaveGrpcMetrics(data *pb.SensorData, md metadata.MD) {

    pSize, mTime := getProtoMetrics(data)
    hSize := 5 + calculateGRPCOverhead(md) 
    lat := calculateLatency(data.Timestamp)

    metricsChan <- Metric{
        Protocol:     "gRPC",
        LatencyMs:    lat,
        PayloadByte:  pSize,
        OverheadByte: hSize,
        MarshalTime:  mTime,
        Timestamp:    time.Now().Format("15:04:05.000"),
    }
}

func resetStats() {
    // Svuota il canale
    for len(metricsChan) > 0 {
        <-metricsChan
    }

    metricsMu.Lock()
    defer metricsMu.Unlock()
    
    warmupUntil = time.Now().Add(warmupDuration)
    history = []Metric{}
    totalPayloadRest = 0
    totalOverheadRest = 0
    totalPayloadGrpc = 0
    totalOverheadGrpc = 0
    atomic.StoreUint64(&msgCountRest, 0)
    atomic.StoreUint64(&msgCountGrpc, 0)
    throughputRest = 0
    throughputGrpc = 0
}