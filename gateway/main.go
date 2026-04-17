package main

import (
    "log"
    "net"
    "net/http"
    "google.golang.org/grpc"
    "time"
    "sync/atomic"
    pb "telemetry-bench/proto"
)

func main() {
    go metricsWorker()
    startThroughputTicker()
    // 1. Avvio gRPC Server
    go func() {
        lis, err := net.Listen("tcp", ":50051")
        if err != nil {
            log.Fatalf("failed to listen: %v", err)
        }
        s := grpc.NewServer()
        pb.RegisterTelemetryServiceServer(s, &telemetryServer{})
        log.Println("gRPC Server in ascolto su :50051")
        s.Serve(lis)
    }()

    // 2. Setup Rotte HTTP
    mux := http.NewServeMux()
    mux.HandleFunc("/results", handleResults)
    mux.HandleFunc("/telemetry", handleTelemetry)
    mux.HandleFunc("/set-mode", handleSetMode)
	mux.HandleFunc("/get-mode", handleGetMode)
	mux.HandleFunc("/set-size", handleSetSize)
	mux.HandleFunc("/get-size", handleGetSize)
    mux.HandleFunc("/set-sensors", handleSetSensors)
	mux.HandleFunc("/get-sensors", handleGetSensors)
    mux.HandleFunc("/set-protocol", handleSetProtocol)
	mux.HandleFunc("/get-protocol", handleGetProtocol)
    mux.HandleFunc("/reset", handleReset)
	mux.HandleFunc("/ws", handleWS)    
	
    fs := http.FileServer(http.Dir("dashboard"))
    mux.Handle("/", fs)

    log.Println("gRPC Server in ascolto su :50051")
    log.Fatal(http.ListenAndServe(":8080", enableCORS(mux)))
}

func startThroughputTicker() {
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        for range ticker.C {
            // Legge e resetta i contatori atomici definiti nelle variabili globali
            currRest := atomic.SwapUint64(&msgCountRest, 0)
            currGrpc := atomic.SwapUint64(&msgCountGrpc, 0)

            metricsMu.Lock()
            throughputRest = float64(currRest)
            throughputGrpc = float64(currGrpc)
            metricsMu.Unlock()
        }
    }()
}