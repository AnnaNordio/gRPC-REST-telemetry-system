package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "telemetry-bench/proto"
)

var (
	mu             sync.Mutex
	stopChannels   = make(map[int]chan struct{})
	currentSensors = 0
	sensorID       = 0
)

func main() {
	log.Println("Avvio Sensore High-Precision Multi-Node...")

	// 1. Setup gRPC Connection (Shared by all sensors)
	conn, err := grpc.Dial(gatewayGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Errore connessione gRPC: %v", err)
	}
	defer conn.Close()
	grpcClient := pb.NewTelemetryServiceClient(conn)

	// 2. HTTP Client Ottimizzato (Shared)
	httpClient := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        500,
			MaxIdleConnsPerHost: 100,
		},
	}

	// Loop di controllo configurazione (Gira ogni secondo)
	for {
		// Recupera la config dal gateway (polling del backend)
		mode, size, targetSensorsStr := fetchConfig(httpClient)

		targetSensors, err := strconv.Atoi(targetSensorsStr)
		if err != nil {
			log.Printf("Errore conversione sensori: %v", err)
			continue
		}

		syncSensors(targetSensors, grpcClient, httpClient, mode, size)

		time.Sleep(1 * time.Second)
	}
}

// syncSensors aggiunge o rimuove goroutine per pareggiare il numero desiderato
func syncSensors(target int, client pb.TelemetryServiceClient, http *http.Client, mode, size string) {
	mu.Lock()
	defer mu.Unlock()

	if target == currentSensors {
		return
	}

	if target > currentSensors {
		// Aggiungi nuovi sensori
		diff := target - currentSensors
		for i := 0; i < diff; i++ {
			sensorID++
			stopCh := make(chan struct{})
			stopChannels[sensorID] = stopCh
			go runVirtualSensor(sensorID, stopCh, client, http, mode, size)
		}
		log.Printf("Scalato a %d sensori", target)
	} else {
		diff := currentSensors - target
		for i := 0; i < diff; i++ {
			for id, ch := range stopChannels {
				close(ch)
				delete(stopChannels, id)
				break
			}
		}
		log.Printf("Scalato a %d sensori", target)
	}
	currentSensors = target
}

// runVirtualSensor rappresenta il ciclo di vita di un singolo sensore virtuale
func runVirtualSensor(id int, stopCh chan struct{}, grpcClient pb.TelemetryServiceClient, httpClient *http.Client, mode, size string) {
	ticker := time.NewTicker(100 * time.Millisecond) // 10Hz
	defer ticker.Stop()

	// In gRPC streaming, ogni sensore apre il suo stream dedicato
	var stream pb.TelemetryService_StreamDataClient

	for {
		select {
		case <-stopCh:
			log.Printf("Sensore [%d] arrestato", id)
			return
		case <-ticker.C:
			data := generateData(size)

			if mode == "polling" {
				// Esecuzione Unary (1Hz approx)
				if time.Now().UnixMilli()%1000 < 100 {
					executePolling(httpClient, grpcClient, data)
				}
			} else {
				// Esecuzione Streaming (10Hz)
				if stream == nil {
					var err error
					stream, err = grpcClient.StreamData(context.Background())
					if err != nil {
						log.Printf("![%d] Errore apertura stream: %v", id, err)
						continue
					}
				}

				if err := stream.Send(data); err != nil {
					log.Printf("![%d] Errore invio stream: %v", id, err)
					stream = nil // Reset per riconnessione
				}

				// REST asincrono per monitorare l'overhead simultaneo
				go sendRest(httpClient, data)
			}
		}
	}
}