package main

import (
    "log"
    "net"
    "net/http"
    "google.golang.org/grpc"
    pb "telemetry-bench/proto"
)

func main() {
    // 1. Avvio gRPC Server
    go func() {
        lis, err := net.Listen("tcp", ":50051")
        if err != nil {
            log.Fatalf("failed to listen: %v", err)
        }
        s := grpc.NewServer()
        pb.RegisterTelemetryServiceServer(s, &telemetryServer{})
        log.Println("🚀 gRPC Server in ascolto su :50051")
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
	mux.HandleFunc("/ws", handleWS)    
	
    fs := http.FileServer(http.Dir("dashboard"))
    mux.Handle("/", fs)

    log.Println("🚀 Gateway e Dashboard attivi su http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", enableCORS(mux)))
}