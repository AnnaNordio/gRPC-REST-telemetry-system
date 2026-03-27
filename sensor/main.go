package main

import (
    "context"
    "log"
    "net/http"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "telemetry-bench/proto"
)

func main() {
    log.Println("🚀 Avvio Sensore High-Precision...")

    // 1. Setup gRPC
    conn, err := grpc.Dial(gatewayGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Errore connessione gRPC: %v", err)
    }
    defer conn.Close()

    grpcClient := pb.NewTelemetryServiceClient(conn)
    stream, _ := grpcClient.StreamData(context.Background())

    // 2. HTTP Client Ottimizzato
    httpClient := &http.Client{
        Timeout: 2 * time.Second,
        Transport: &http.Transport{MaxIdleConns: 100, MaxIdleConnsPerHost: 100},
    }

    ticker := time.NewTicker(100 * time.Millisecond)
    var currentMode, currentSize = "polling", "small"
    lastConfigCheck := time.Now()

    for range ticker.C {
        // Refresh configurazione ogni secondo
        if time.Since(lastConfigCheck) > time.Second {
            currentMode, currentSize = fetchConfig(httpClient)
            lastConfigCheck = time.Now()
        }

        data := generateData(currentSize)

        if currentMode == "polling" {
            // Polling: invio ogni secondo circa (frequenza 1Hz)
            if time.Now().UnixMilli()%1000 < 100 {
                executePolling(httpClient, grpcClient, data)
            }
        } else {
            // Streaming: invio continuo (frequenza 10Hz)
            executeStreaming(httpClient, stream, data)
        }
    }
}