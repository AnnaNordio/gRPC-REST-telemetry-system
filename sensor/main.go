package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	log.Println("🚀 Avvio Sensore Modulare...")

	// 1. Setup gRPC (Connessione)
	conn, err := grpc.Dial(gatewayGrpcAddr, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("Impossibile connettersi a gRPC: %v", err)
	}
	defer conn.Close()
	grpcClient := pb.NewTelemetryServiceClient(conn)

	// 2. Setup HTTP Client
	httpClient := &http.Client{Timeout: 2 * time.Second}

	for {
		data := generateData()

		// --- ESECUZIONE REST ---
		sizeREST := getJsonSize(data)
		startR := time.Now()
		sendRest(httpClient, data)
		latR := float64(time.Since(startR).Microseconds())
		updateLatency("REST", latR) 

		// --- ESECUZIONE gRPC ---
		sizeGRPC := getProtoSize(data)
		startG := time.Now()
		_, err := grpcClient.SendData(context.Background(), data)
		latG := float64(time.Since(startG).Microseconds())
		
		if err != nil {
			log.Printf("Errore gRPC: %v", err)
		} else {
			updateLatency("gRPC", latG) 
			fmt.Printf("[%s] REST: %.0fµs (%0.fB) | gRPC: %.0fµs (%0.fB)\n", 
				time.Now().Format("15:04:05"), latR, sizeREST, latG, sizeGRPC)
		}

		time.Sleep(1 * time.Second)
	}
}

// Funzione di supporto per l'invio REST
func sendRest(client *http.Client, data *pb.SensorData) {
	b, _ := json.Marshal(data)
	
	// Se hai creato payload.go, qui potresti chiamare:
	// savePayload("REST", float64(len(b)))

	req, _ := http.NewRequest("POST", gatewayRestAddr, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}
}