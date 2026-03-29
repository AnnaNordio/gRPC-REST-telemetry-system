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

// Stato globale della configurazione per i sensori
var (
	configMu sync.RWMutex
	globalMode string
	globalSize string

	sensorMu       sync.Mutex
	stopChannels   = make(map[int]chan struct{})
	currentSensors = 0
	sensorID       = 0
)

func main() {
	log.Println("Avvio Sensore High-Precision Multi-Node...")

	// 1. Setup gRPC Connection (Condivisa tra tutti i sensori)
	conn, err := grpc.Dial(gatewayGrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Errore connessione gRPC: %v", err)
	}
	defer conn.Close()
	grpcClient := pb.NewTelemetryServiceClient(conn)

	// 2. HTTP Client Ottimizzato (Condiviso)
	httpClient := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        500,
			MaxIdleConnsPerHost: 100,
		},
	}

	// Loop di controllo configurazione (Polling ogni secondo)
    for {
        // 1. Recupera l'oggetto configurazione intero
        config := fetchFullConfig(httpClient)

        // 2. AGGIORNAMENTO ATOMICO della configurazione globale
        // Usiamo i campi della struct (config.Mode e config.Size)
        configMu.Lock()
        globalMode = config.Mode
        globalSize = config.Size
        configMu.Unlock()

        log.Printf("Config aggiornata: mode=%s, size=%s, sensors=%d", 
            config.Mode, config.Size, config.Sensors)

        // 3. Sincronizza il numero di goroutine usando config.Sensors
        // Nota: Sensors è già un int, non serve più strconv.Atoi qui!
        syncSensors(config.Sensors, grpcClient, httpClient)

        time.Sleep(1 * time.Second)
    }
}

func syncSensors(target int, client pb.TelemetryServiceClient, http *http.Client) {
	sensorMu.Lock()
	defer sensorMu.Unlock()

	if target == currentSensors {
		return
	}

	if target > currentSensors {
		diff := target - currentSensors
		for i := 0; i < diff; i++ {
			sensorID++
			stopCh := make(chan struct{})
			stopChannels[sensorID] = stopCh
			// Ora passiamo solo i client, mode e size verranno letti dinamicamente
			go runVirtualSensor(sensorID, stopCh, client, http)
		}
	} else {
		diff := currentSensors - target
		for i := 0; i < diff; i++ {
			for id, ch := range stopChannels {
				close(ch)
				delete(stopChannels, id)
				break
			}
		}
	}
	currentSensors = target
	log.Printf("Flotta sensori aggiornata: %d attivi", currentSensors)
}

func runVirtualSensor(id int, stopCh chan struct{}, grpcClient pb.TelemetryServiceClient, httpClient *http.Client) {
	ticker := time.NewTicker(100 * time.Millisecond) // 10Hz
	defer ticker.Stop()

	var stream pb.TelemetryService_StreamDataClient
	var lastMode string

	for {
		select {
		case <-stopCh:
			if stream != nil {
				stream.CloseSend()
			}
			log.Printf("Sensore [%d] arrestato", id)
			return
		case <-ticker.C:
			// LETTURA DINAMICA della configurazione (RLock per massime performance)
			configMu.RLock()
			mode := globalMode
			size := globalSize
			configMu.RUnlock()

			// Se il modo cambia, resettiamo lo stream gRPC se esistente
			if mode != lastMode && stream != nil {
				stream.CloseSend()
				stream = nil
			}
			lastMode = mode
			log.Printf("Sensore [%d] generazione dati: mode=%s, size=%s", id, mode, size)
			data := generateData(size)

			if mode == "polling" {
				// Esecuzione Unary (REQ-RES)
				// Eseguiamo solo 1 volta al secondo per non saturare i log in Unary
				if time.Now().UnixMilli()%1000 < 100 {
					executePolling(httpClient, grpcClient, data)
				}
			} else {
				// Esecuzione STREAMING
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
					stream = nil // Reset per riconnessione al prossimo tick
				}

				// REST simultaneo per confronto overhead (come nel tuo codice originale)
				go sendRest(httpClient, data)
			}
		}
	}
}