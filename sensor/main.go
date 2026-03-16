package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "telemetry-bench/proto" 
)

const (
	gatewayRestAddr = "http://localhost:8080/telemetry"
	gatewayGrpcAddr = "localhost:50051"
)

func main() {
	log.Println("--- Avvio Sensore ---")

	// 1. Connessione gRPC (con timeout per non restare appesi)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, gatewayGrpcAddr, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), // Attende che la connessione sia stabilita
	)
	if err != nil {
		log.Fatalf("ERRORE: Gateway gRPC non raggiungibile: %v", err)
	}
	defer conn.Close()
	grpcClient := pb.NewTelemetryServiceClient(conn)
	log.Println("Connesso a gRPC Gateway")

	// 2. Client HTTP
	httpClient := &http.Client{Timeout: 2 * time.Second}
	log.Println("Client HTTP pronto")

	// 3. Loop di invio
	log.Println("Inizio invio dati ogni secondo...")
	for {
		data := &pb.SensorData{
			SensorId:    "sensor-01",
			Temperature: 20 + rand.Float32()*10,
			Humidity:    40 + rand.Float32()*20,
			Timestamp:   time.Now().UnixMilli(),
		}

		// Test REST
		startR := time.Now()
		sendRest(httpClient, data)
		latR := time.Since(startR)

		// Test gRPC
		startG := time.Now()
		_, err := grpcClient.SendData(context.Background(), data)
		latG := time.Since(startG)

		if err != nil {
			log.Printf("Errore invio gRPC: %v", err)
		} else {
			fmt.Printf("[%s] SUCCESS | REST: %v | gRPC: %v\n", 
				time.Now().Format("15:04:05"), latR, latG)
		}

		time.Sleep(1 * time.Second)
	}
}

func sendRest(client *http.Client, data *pb.SensorData) {
	b, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", gatewayRestAddr, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Errore REST: %v", err)
		return
	}
	defer resp.Body.Close()
}