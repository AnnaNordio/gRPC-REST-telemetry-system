package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"telemetry-bench/pkg/config"
	pb "telemetry-bench/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	globalConfig atomic.Value

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

	//Setup gRPC Connection Pool
	const poolSize = 100
	var grpcClients []pb.TelemetryServiceClient

	for i := 0; i < poolSize; i++ {
		conn, err := grpc.NewClient(
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

	// HTTP
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
		return
	}

	for {
		config := fetchFullConfig(httpClient)

		globalConfig.Store(config)

		syncSensors(config.Sensors, grpcClients, httpClient)

		time.Sleep(2 * time.Second)
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

			selectedClient := clients[sensorID%len(clients)]
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
			conf := globalConfig.Load().(config.TelemetryConfig)
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
