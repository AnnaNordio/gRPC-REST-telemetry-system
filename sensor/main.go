package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "telemetry-bench/proto"
)

const (
	gatewayRestAddr = "http://localhost:8080/telemetry"
	gatewayGrpcAddr = "localhost:50051"
	modeEndpoint    = "http://localhost:8080/get-mode"
)

func main() {
	log.Println("🚀 Avvio Sensore Modulare (Polling/Streaming)...")

	// 1. Setup gRPC
	conn, err := grpc.Dial(gatewayGrpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Impossibile connettersi a gRPC: %v", err)
	}
	defer conn.Close()
	grpcClient := pb.NewTelemetryServiceClient(conn)

	// Creazione dello stream gRPC (per modalità streaming)
	stream, err := grpcClient.StreamData(context.Background())
	if err != nil {
		log.Fatalf("Errore apertura stream gRPC: %v", err)
	}

	// 2. Setup HTTP Client Ottimizzato per Keep-Alive
	httpClient := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			MaxIdleConnsPerHost: 100,
		},
	}

	// 3. Loop di controllo e invio
	ticker := time.NewTicker(10 * time.Millisecond) // Risoluzione 10ms
	defer ticker.Stop()

	var currentMode string = "polling"
	lastModeCheck := time.Now()

	for range ticker.C {
		// Controllo la modalità dal Gateway ogni secondo
		if time.Since(lastModeCheck) > time.Second {
			currentMode = fetchMode(httpClient)
			lastModeCheck = time.Now()
		}

		data := generateData() // Funzione che genera i tuoi dati

		if currentMode == "polling" {
			// In POLLING inviamo una volta al secondo (ogni 100 cicli da 10ms)
			if time.Now().UnixNano()/int64(time.Millisecond)%1000 < 10 {
				executePolling(httpClient, grpcClient, data)
			}
		} else {
			// In STREAMING spariamo a raffica (ogni 10ms)
			executeStreaming(httpClient, stream, data)
		}
	}
}

// --- LOGICA POLLING (Unary) ---
func executePolling(client *http.Client, grpcClient pb.TelemetryServiceClient, data *pb.SensorData) {
	// REST
	fmt.Printf("[POLLING]")
	startR := time.Now()
	sendRest(client, data)
	latR := float64(time.Since(startR).Microseconds())
	updateLatency("REST", latR)
	fmt.Printf("REST: %7.0f µs\n", latR)

	// gRPC Unary
	startG := time.Now()
	_, err := grpcClient.SendData(context.Background(), data)
	if err == nil {
		latG := float64(time.Since(startG).Microseconds())
		updateLatency("gRPC", latG)
		fmt.Printf("gRPC: %7.0f µs\n", latG)
	}
	
}

// --- LOGICA STREAMING ---
func executeStreaming(client *http.Client, stream pb.TelemetryService_StreamDataClient, data *pb.SensorData) {
	// REST Streaming (Simulato con goroutine asincrona)
	fmt.Printf("[STREAMING]")
	go func(d *pb.SensorData) {
		startR := time.Now()
		sendRest(client, d)
		latR := float64(time.Since(startR).Microseconds())
		updateLatency("REST", latR)
		fmt.Printf("REST: %7.0f µs\n", latR)
	}(data)

	// gRPC Real Streaming (Sullo stream aperto)
	startG := time.Now()
	err := stream.Send(data)
	if err == nil {
		latG := float64(time.Since(startG).Microseconds())
		updateLatency("gRPC", latG)
		fmt.Printf("gRPC: %7.0f µs\n", latG)
	}
}


// Helper per recuperare la modalità dal gateway
func fetchMode(client *http.Client) string {
	resp, err := client.Get(modeEndpoint)
	if err != nil {
		return "polling"
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

// Helper invio REST
func sendRest(client *http.Client, data *pb.SensorData) {
	b, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", gatewayRestAddr, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err == nil {
		// Importante: svuotare e chiudere il body per riutilizzare la connessione TCP
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}