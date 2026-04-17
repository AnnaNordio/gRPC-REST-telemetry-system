package main

import (
    "time"
    "net/http"
    "google.golang.org/grpc/metadata"  
    "sync/atomic"
    "os"
    pb "telemetry-bench/proto"
)

// Canale per trasportare le metriche dai server al worker
// Capacità 10.000 per gestire picchi di 100 sensori a 10Hz
var metricsChan = make(chan Metric, 10000)

func metricsWorker() {
    os.MkdirAll("results", 0755)
    
    writer := &MetricsWriter{}

    flushTicker := time.NewTicker(1 * time.Second)

    for {
        select {
        case m, ok := <-metricsChan:
            if !ok { return }
            
            // --- 1. Lettura Stato (dal tuo file state) ---
            metricsMu.Lock()
            isWarmup := time.Now().Before(warmupUntil)
            mMode := currentMode
            mSize := currentSize
            mProto := currentProtocol
            mSensors := currentSensors
            metricsMu.Unlock()

            if isWarmup {
                continue
            }

            // --- 2. Aggiornamento Statistiche (Usa le tue variabili di state) ---
            updateStats(m)

            // --- 3. Scrittura su File ---
            writer.Write(m, mMode, mSize, mProto, mSensors)

        case <-flushTicker.C:
            if writer.csvWriter != nil {
                writer.csvWriter.Flush() // Scrive tutto il blocco accumulato in una volta sola
            }
        }
        
    }
}

func updateStats(m Metric) {
    metricsMu.Lock()
    defer metricsMu.Unlock()

    if m.Protocol == "gRPC" {
        atomic.AddUint64(&msgCountGrpc, 1)
        totalPayloadGrpc += m.PayloadByte
        totalOverheadGrpc += m.OverheadByte
    } else {
        atomic.AddUint64(&msgCountRest, 1)
        totalPayloadRest += m.PayloadByte
        totalOverheadRest += m.OverheadByte
    }

    // Aggiunta alla history per i grafici live
    history = append(history, m)
    
    // Cleanup history vecchia (es. mantieni ultimi 1000 punti)
    if len(history) > 1000 {
        history = history[len(history)-1000:]
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