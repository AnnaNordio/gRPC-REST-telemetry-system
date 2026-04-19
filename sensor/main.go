package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
    "os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "telemetry-bench/proto"
	"telemetry-bench/pkg/config"

)

// Variabili globali per il coordinamento
var (
	globalConfig atomic.Value // Contiene l'ultima TelemetryConfig caricata

	sensorMu       sync.Mutex
	stopChannels   = make(map[int]chan struct{})
	currentSensors = 0
	sensorID       = 0
)

func main() {
	log.Println("Avvio Sensore High-Precision Multi-Node...")

	globalConfig.Store(config.TelemetryConfig{
		Mode:     "polling",
		Size:     "small",
		Sensors:  0,
		Protocol: "both",
	})

	// 1. Setup gRPC Connection Pool (100 connessioni TCP separate)
	const poolSize = 100
	var grpcClients []pb.TelemetryServiceClient

	for i := 0; i < poolSize; i++ {
		conn, err := grpc.Dial(
			gatewayGrpcAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithInitialWindowSize(1<<20),
			grpc.WithInitialConnWindowSize(1<<20),
		)
		if err != nil {
			log.Fatalf("Errore connessione gRPC pool: %v", err)
		}
		grpcClients = append(grpcClients, pb.NewTelemetryServiceClient(conn))
	}

	// 2. HTTP Client Ottimizzato per polling e fetch
	httpClient := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        500,
			MaxIdleConnsPerHost: 100,
		},
	}

    appMode := os.Getenv("APP_MODE")

    if appMode == "benchmark" {
        log.Println(">>> MODALITÀ BENCHMARK RILEVATA <<<")
        runBenchmarkSuite(grpcClients, httpClient)
        log.Println("Benchmark completato. Spegnimento in corso...")
        time.Sleep(2 * time.Second)
        return // Esci dal programma, Docker fermerà il container
    }

	// 3. Loop di monitoraggio configurazione
	for {
		config := fetchFullConfig(httpClient)

		// AGGIORNAMENTO ATOMICO: i sensori leggono questo valore senza lock
		globalConfig.Store(config)

		// Gestione del numero di goroutine attive
		syncSensors(config.Sensors, grpcClients, httpClient)

		time.Sleep(2 * time.Second) // Controlla la config ogni 2s per non sovraccaricare
	}
}

func syncSensors(target int, clients []pb.TelemetryServiceClient, httpClient *http.Client) {
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
			
			// Distribuzione Round-Robin sul pool gRPC
			selectedClient := clients[sensorID % len(clients)]
			go runVirtualSensor(sensorID, stopCh, selectedClient, httpClient)
		}
		log.Printf("[Sync] Attivati %d nuovi sensori. Totale: %d", diff, target)
	} else {
		diff := currentSensors - target
		for i := 0; i < diff; i++ {
			for id, ch := range stopChannels {
				close(ch)
				delete(stopChannels, id)
				break
			}
		}
		log.Printf("[Sync] Rimossi %d sensori. Totale: %d", diff, target)
	}
	currentSensors = target
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
			return
		case <-ticker.C:
			// Lettura lock-free della configurazione globale
			conf := globalConfig.Load().(config.TelemetryConfig)
			// Se la modalità cambia, dobbiamo resettare lo stream gRPC
			if conf.Mode != lastMode && stream != nil {
				stream.CloseSend()
				stream = nil
			}
			lastMode = conf.Mode

			data := generateData(conf.Size)

			if conf.Mode == "polling" {
				executePolling(httpClient, grpcClient, data, conf.Protocol)
			} else {
				if stream == nil {
					var err error
					stream, err = grpcClient.StreamData(context.Background())
					if err != nil {
						log.Printf("[%d] Errore Stream: %v", id, err)
						continue
					}
				}
				executeStreaming(httpClient, stream, data, conf.Protocol)
			}
		}
	}
}