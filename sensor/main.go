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
    globalProtocol string

	sensorMu       sync.Mutex
	stopChannels   = make(map[int]chan struct{})
	currentSensors = 0
	sensorID       = 0
)

func main() {
    log.Println("Avvio Sensore High-Precision Multi-Node...")

    // 1. Setup gRPC Connection Pool (8 connessioni TCP separate)
    const poolSize = 100
    var grpcClients []pb.TelemetryServiceClient
    
    for i := 0; i < poolSize; i++ {
        conn, err := grpc.Dial(
            gatewayGrpcAddr, 
            grpc.WithTransportCredentials(insecure.NewCredentials()),
            // Aumentiamo i buffer per gestire i dati "large" senza strozzature
            grpc.WithInitialWindowSize(1 << 20),     
            grpc.WithInitialConnWindowSize(1 << 20), 
        )
        if err != nil {
            log.Fatalf("Errore connessione gRPC pool: %v", err)
        }
        // NOTA: In un'app reale dovresti gestire la chiusura di tutte le conn nel pool
        grpcClients = append(grpcClients, pb.NewTelemetryServiceClient(conn))
    }

    // 2. HTTP Client Ottimizzato (Condiviso)
    httpClient := &http.Client{
        Timeout: 2 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:        500,
            MaxIdleConnsPerHost: 100, // REST usa fino a 100 socket paralleli
        },
    }

    go func() {
        for {
            time.Sleep(5 * time.Second)

            sensorMu.Lock()
            totalSensors := currentSensors
            sensorMu.Unlock()

            log.Printf("--- STATS CONNESSIONI ---")
            log.Printf("Sensori Attivi: %d", totalSensors)
            
            // Per gRPC, dato il tuo round-robin:
            usedGrpc := totalSensors
            if usedGrpc > poolSize {
                usedGrpc = poolSize
            }
            log.Printf("[gRPC] Connessioni nel pool utilizzate: %d/%d", usedGrpc, poolSize)

            // Per REST, leggiamo le idle connections dal transport
            if transport, ok := httpClient.Transport.(*http.Transport); ok {
                // Nota: Go non espone facilmente le connessioni "attive", 
                // ma quelle "Idle" (aperte e pronte al riuso) sono un ottimo indicatore.
                log.Printf("[REST] Connessioni in stato Idle (pronte): %d", transport.MaxIdleConnsPerHost)
            }
            log.Printf("--------------------------")
        }
    }()

    for {
        config := fetchFullConfig(httpClient)

        configMu.Lock()
        globalMode = config.Mode
        globalSize = config.Size
        globalProtocol = config.Protocol
        configMu.Unlock()

        // Passiamo l'intero POOL di client a syncSensors
        syncSensors(config.Sensors, grpcClients, httpClient)

        time.Sleep(1 * time.Second)
    }
}

// Aggiornata la firma per accettare la slice di client
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
            
            // ASSEGNAZIONE ROUND-ROBIN: 
            // Distribuiamo i 100 sensori sulle 8 connessioni gRPC
            selectedClient := clients[sensorID % len(clients)]
            
            go runVirtualSensor(sensorID, stopCh, selectedClient, httpClient)
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
}

func runVirtualSensor(id int, stopCh chan struct{}, grpcClient pb.TelemetryServiceClient, httpClient *http.Client) {
    // Ticker a 100ms = 10Hz (10 messaggi al secondo)
    ticker := time.NewTicker(100 * time.Millisecond)
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
            // 1. Lettura configurazione dinamica
            configMu.RLock()
            mode := globalMode
            size := globalSize
            configMu.RUnlock()

            // 2. Gestione cambio modalità e reset stream
            if mode != lastMode && stream != nil {
                stream.CloseSend()
                stream = nil
            }
            lastMode = mode

            // 3. Generazione dati
            data := generateData(size)

            // 4. Esecuzione Logica di Invio
            if mode == "polling" {
                // Esecuzione Unary (REQ-RES) a 10Hz
                // Rimosso il filtro %1000 per uniformare la frequenza
                executePolling(httpClient, grpcClient, data, globalProtocol)
            } else {
                // Esecuzione STREAMING a 10Hz
                // Inizializzazione stream se necessario
                if stream == nil {
                    var err error
                    stream, err = grpcClient.StreamData(context.Background())
                    if err != nil {
                        log.Printf("![%d] Errore apertura stream: %v", id, err)
                        continue
                    }
                }

                executeStreaming(httpClient, stream, data, globalProtocol)
            }
        }
    }
}