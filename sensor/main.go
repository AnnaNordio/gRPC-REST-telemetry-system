package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "telemetry-bench/proto"
)

const (
	gatewayRestAddr = "http://gateway:8080/telemetry"
	gatewayGrpcAddr = "gateway:50051"
	modeEndpoint    = "http://gateway:8080/get-mode"
)

func main() {
	log.Println("🚀 Avvio Sensore High-Precision (Timestamp Sync)...")

	// 1. Setup gRPC
	conn, err := grpc.Dial(gatewayGrpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Impossibile connettersi a gRPC: %v", err)
	}
	defer conn.Close()
	grpcClient := pb.NewTelemetryServiceClient(conn)

	stream, err := grpcClient.StreamData(context.Background())
	if err != nil {
		log.Fatalf("Errore apertura stream gRPC: %v", err)
	}

	// 2. HTTP Client Ottimizzato
	httpClient := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
		},
	}

	// 3. Loop di controllo (100ms per benchmark fluido)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	var currentMode string = "polling"
	lastModeCheck := time.Now()

	for range ticker.C {
		if time.Since(lastModeCheck) > time.Second {
			currentMode = fetchMode(httpClient)
			lastModeCheck = time.Now()
		}

		data := generateData()

		if currentMode == "polling" {
			// Polling: un invio ogni secondo
			if time.Now().UnixMilli()%1000 < 100 {
				executePolling(httpClient, grpcClient, data)
			}
		} else {
			// Streaming: invio continuo
			executeStreaming(httpClient, stream, data)
		}
	}
}

// --- LOGICA STREAMING (Precisione Microsecondi) ---
func executeStreaming(client *http.Client, stream pb.TelemetryService_StreamDataClient, data *pb.SensorData) {
	startTs := time.Now().UnixMicro()
	data.Timestamp = startTs

	if err := stream.Send(data); err != nil {
		fmt.Printf("gRPC Send Error: %v\n", err)
	}
	go func(d pb.SensorData) {
		sendRest(client, &d)
	}(*data)

	fmt.Printf("[SENT] ID: %s (gRPC + REST)\n", data.Timestamp)
}

// --- LOGICA POLLING (Unary) ---
func executePolling(client *http.Client, grpcClient pb.TelemetryServiceClient, data *pb.SensorData) {
	ts := time.Now().UnixMicro()
	data.Timestamp = ts

	// Invio asincrono di entrambi per non bloccare il loop
	go sendRest(client, data)
	go func() {
		_, _ = grpcClient.SendData(context.Background(), data)
	}()
	
	fmt.Printf("[POLLING SENT] ID: %s\n", data.Timestamp)
}

// --- HELPERS ---

func sendRest(client *http.Client, data *pb.SensorData) {
	b, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", gatewayRestAddr, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err == nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func fetchMode(client *http.Client) string {
	resp, err := client.Get(modeEndpoint)
	if err != nil {
		return "polling"
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}