package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"sync/atomic"
	pb "telemetry-bench/proto"
	"time"

	"google.golang.org/grpc"
)

func main() {
	go metricsWorker()
	startThroughputTicker()
	//gRPC Server
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
	isBenchMode := os.Getenv("BENCH_MODE") == "true"
	// HTTP Server
	mux := http.NewServeMux()

	mux.HandleFunc("/results", handleResults)
	mux.HandleFunc("/get-mode", handleGetMode)
	mux.HandleFunc("/get-size", handleGetSize)
	mux.HandleFunc("/get-sensors", handleGetSensors)
	mux.HandleFunc("/get-protocol", handleGetProtocol)
	mux.HandleFunc("/get-config", handleGetConfig)
	mux.HandleFunc("/set-config", handleSetConfig)
	mux.HandleFunc("/telemetry", handleTelemetry)

	if !isBenchMode {
		// --- MODALITÀ INTERATTIVA ---
		mux.HandleFunc("/set-mode", handleSetMode)
		mux.HandleFunc("/set-size", handleSetSize)
		mux.HandleFunc("/set-sensors", handleSetSensors)
		mux.HandleFunc("/set-protocol", handleSetProtocol)
		mux.HandleFunc("/reset", handleReset)
		mux.HandleFunc("/ws", handleWS)

		fs := http.FileServer(http.Dir("dashboard"))
		mux.Handle("/", fs)
	} else {
		// --- MODALITÀ BENCHMARK ---
		log.Println("BENCHMARK MODE")
	}

	log.Println("HTTP Server listening on :8080")
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
