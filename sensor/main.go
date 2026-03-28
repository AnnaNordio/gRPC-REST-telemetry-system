package main

import (
	"context"
	"log"
	"net/http"
	"sync"
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
	
	// Gestiamo lo stream in modo che possa essere ricreato se cade
	var stream pb.TelemetryService_StreamDataClient
	var streamMu sync.Mutex

	// 2. HTTP Client Ottimizzato
	httpClient := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	var currentMode, currentSize = "polling", "small"
	lastConfigCheck := time.Now()

	for range ticker.C {
		// Refresh configurazione
		if time.Since(lastConfigCheck) > time.Second {
			currentMode, currentSize = fetchConfig(httpClient)
			lastConfigCheck = time.Now()
		}

		// IMPORTANTE: Generiamo i dati freschi per ogni ciclo
		data := generateData(currentSize)

		if currentMode == "polling" {
			// Invio Unary (Simula Polling) - circa 1Hz
			if time.Now().UnixMilli()%1000 < 100 {
				// Passiamo una copia o dati freschi per evitare race conditions nelle goroutine
				go executePolling(httpClient, grpcClient, generateData(currentSize))
			}
		} else {
			// Invio Streaming - 10Hz
			streamMu.Lock()
			if stream == nil {
				stream, _ = grpcClient.StreamData(context.Background())
			}
			
			// Se Send fallisce, resettiamo lo stream per il prossimo tentativo
			if err := stream.Send(data); err != nil {
				log.Printf("Stream error: %v. Reconnecting...", err)
				stream = nil 
			}
			streamMu.Unlock()

			// REST rimane asincrono per non rallentare il loop di streaming
			go sendRest(httpClient, data)
		}
	}
}