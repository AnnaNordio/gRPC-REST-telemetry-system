package main

import (
    "time"
    "net/http"
    "google.golang.org/grpc/metadata"  
    "sync/atomic"
    "encoding/csv"
    "os"
    "log"
    "strconv"
    pb "telemetry-bench/proto"
)

// Canale per trasportare le metriche dai server al worker
// Capacità 10.000 per gestire picchi di 100 sensori a 10Hz
var metricsChan = make(chan Metric, 10000)

// Il Worker: l'unico punto che scrive nella history, eliminando la contesa del Mutex
func metricsWorker() {
    // 1. APRI IL FILE (Fuori dal loop)
    os.MkdirAll("results", 0755)
    log.Println("Worker delle metriche avviato, pronto a ricevere dati...")
    file, err := os.OpenFile("results/bench_results.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    
    if err != nil {
        return 
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // 2. SCRIVI L'HEADER SOLO SE IL FILE È VUOTO (Fuori dal loop)
    info, _ := file.Stat()
    if info.Size() == 0 {
        writer.Write([]string{"Timestamp", "Protocol", "LatencyMs", "PayloadBytes", "OverheadBytes", "MarshalTimeMs"})
        writer.Flush() // Scrive l'intestazione su disco immediatamente
    }

    // 3. IL LOOP CHE GESTISCE I DATI
    for m := range metricsChan {
        if time.Now().Before(warmupUntil) {
            continue 
        }

        // --- Logica statistiche (Invariata) ---
        metricsMu.Lock()
        if m.Protocol == "gRPC" {
            atomic.AddUint64(&msgCountGrpc, 1)
            totalPayloadGrpc += m.PayloadByte
            totalOverheadGrpc += m.OverheadByte
        } else {
            atomic.AddUint64(&msgCountRest, 1)
            totalPayloadRest += m.PayloadByte
            totalOverheadRest += m.OverheadByte
        }
        history = append(history, m)
        cutoff := time.Now().Add(-30 * time.Second).Format("15:04:05.000")
        if len(history) > 0 && history[0].Timestamp < cutoff {
            history = history[1:]
        }
        metricsMu.Unlock()

        // --- SCRITTURA DATI (Dentro il loop, scrive solo i valori) ---
        record := []string{
            m.Timestamp,
            m.Protocol,
            strconv.FormatFloat(m.LatencyMs, 'f', 4, 64),
            strconv.FormatInt(m.PayloadByte, 10),
            strconv.FormatInt(m.OverheadByte, 10),
            strconv.FormatFloat(m.MarshalTime, 'f', 6, 64),
        }
        
        writer.Write(record)
        writer.Flush() // Ora scrive solo la riga dei dati
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