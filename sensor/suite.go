package main

import (
	"log"
	"net/http"
	"time"
	"bytes"
	"encoding/json"
	pb "telemetry-bench/proto"
	"telemetry-bench/pkg/config"
)

type TestCase struct {
	Sensors  int
	Mode     string
	Payload  string 
	Protocol string
}

func runBenchmarkSuite(clients []pb.TelemetryServiceClient, httpClient *http.Client) {
	// Qui configuri la tua "tabella di marcia"
	suite := []TestCase{
		{Sensors: 1, Mode: "polling", Payload: "small", Protocol: "grpc"},
		{Sensors: 1, Mode: "polling", Payload: "medium", Protocol: "grpc"},
		{Sensors: 1, Mode: "polling", Payload: "large", Protocol: "grpc"},
		{Sensors: 1, Mode: "polling", Payload: "nested", Protocol: "grpc"},
		{Sensors: 1, Mode: "polling", Payload: "small", Protocol: "rest"},
		{Sensors: 1, Mode: "polling", Payload: "medium", Protocol: "rest"},
		{Sensors: 1, Mode: "polling", Payload: "large", Protocol: "rest"},
		{Sensors: 1, Mode: "polling", Payload: "nested", Protocol: "rest"},

		{Sensors: 10, Mode: "polling", Payload: "small", Protocol: "grpc"},
		{Sensors: 10, Mode: "polling", Payload: "medium", Protocol: "grpc"},
		{Sensors: 10, Mode: "polling", Payload: "large", Protocol: "grpc"},
		{Sensors: 10, Mode: "polling", Payload: "nested", Protocol: "grpc"},
		{Sensors: 10, Mode: "polling", Payload: "small", Protocol: "rest"},
		{Sensors: 10, Mode: "polling", Payload: "medium", Protocol: "rest"},
		{Sensors: 10, Mode: "polling", Payload: "large", Protocol: "rest"},
		{Sensors: 10, Mode: "polling", Payload: "nested", Protocol: "rest"},

		{Sensors: 50, Mode: "polling", Payload: "small", Protocol: "grpc"},
		{Sensors: 50, Mode: "polling", Payload: "medium", Protocol: "grpc"},
		{Sensors: 50, Mode: "polling", Payload: "large", Protocol: "grpc"},
		{Sensors: 50, Mode: "polling", Payload: "nested", Protocol: "grpc"},
		{Sensors: 50, Mode: "polling", Payload: "small", Protocol: "rest"},
		{Sensors: 50, Mode: "polling", Payload: "medium", Protocol: "rest"},
		{Sensors: 50, Mode: "polling", Payload: "large", Protocol: "rest"},
		{Sensors: 50, Mode: "polling", Payload: "nested", Protocol: "rest"},

		{Sensors: 100, Mode: "polling", Payload: "small", Protocol: "grpc"},
		{Sensors: 100, Mode: "polling", Payload: "medium", Protocol: "grpc"},
		{Sensors: 100, Mode: "polling", Payload: "large", Protocol: "grpc"},
		{Sensors: 100, Mode: "polling", Payload: "nested", Protocol: "grpc"},
		{Sensors: 100, Mode: "polling", Payload: "small", Protocol: "rest"},
		{Sensors: 100, Mode: "polling", Payload: "medium", Protocol: "rest"},
		{Sensors: 100, Mode: "polling", Payload: "large", Protocol: "rest"},
		{Sensors: 100, Mode: "polling", Payload: "nested", Protocol: "rest"},


		{Sensors: 1, Mode: "streaming", Payload: "small", Protocol: "grpc"},
		{Sensors: 1, Mode: "streaming", Payload: "medium", Protocol: "grpc"},
		{Sensors: 1, Mode: "streaming", Payload: "large", Protocol: "grpc"},
		{Sensors: 1, Mode: "streaming", Payload: "nested", Protocol: "grpc"},
		{Sensors: 1, Mode: "streaming", Payload: "small", Protocol: "rest"},
		{Sensors: 1, Mode: "streaming", Payload: "medium", Protocol: "rest"},
		{Sensors: 1, Mode: "streaming", Payload: "large", Protocol: "rest"},
		{Sensors: 1, Mode: "streaming", Payload: "nested", Protocol: "rest"},

		{Sensors: 10, Mode: "streaming", Payload: "small", Protocol: "grpc"},
		{Sensors: 10, Mode: "streaming", Payload: "medium", Protocol: "grpc"},
		{Sensors: 10, Mode: "streaming", Payload: "large", Protocol: "grpc"},
		{Sensors: 10, Mode: "streaming", Payload: "nested", Protocol: "grpc"},
		{Sensors: 10, Mode: "streaming", Payload: "small", Protocol: "rest"},
		{Sensors: 10, Mode: "streaming", Payload: "medium", Protocol: "rest"},
		{Sensors: 10, Mode: "streaming", Payload: "large", Protocol: "rest"},
		{Sensors: 10, Mode: "streaming", Payload: "nested", Protocol: "rest"},

		{Sensors: 50, Mode: "streaming", Payload: "small", Protocol: "grpc"},
		{Sensors: 50, Mode: "streaming", Payload: "medium", Protocol: "grpc"},
		{Sensors: 50, Mode: "streaming", Payload: "large", Protocol: "grpc"},
		{Sensors: 50, Mode: "streaming", Payload: "nested", Protocol: "grpc"},
		{Sensors: 50, Mode: "streaming", Payload: "small", Protocol: "rest"},
		{Sensors: 50, Mode: "streaming", Payload: "medium", Protocol: "rest"},
		{Sensors: 50, Mode: "streaming", Payload: "large", Protocol: "rest"},
		{Sensors: 50, Mode: "streaming", Payload: "nested", Protocol: "rest"},

		{Sensors: 100, Mode: "streaming", Payload: "small", Protocol: "grpc"},
		{Sensors: 100, Mode: "streaming", Payload: "medium", Protocol: "grpc"},
		{Sensors: 100, Mode: "streaming", Payload: "large", Protocol: "grpc"},
		{Sensors: 100, Mode: "streaming", Payload: "nested", Protocol: "grpc"},
		{Sensors: 100, Mode: "streaming", Payload: "small", Protocol: "rest"},
		{Sensors: 100, Mode: "streaming", Payload: "medium", Protocol: "rest"},
		{Sensors: 100, Mode: "streaming", Payload: "large", Protocol: "rest"},
		{Sensors: 100, Mode: "streaming", Payload: "nested", Protocol: "rest"},
	}

	for _, tc := range suite {
		log.Printf("\n>>> AVVIO TEST CASE: Sensori:%d | Mode:%s | Size:%s | Protocol:%s <<<", 
            tc.Sensors, tc.Mode, tc.Payload, tc.Protocol)
		updateGatewayConfig(httpClient, tc)
		// 1. Applica la configurazione in modo atomico
		globalConfig.Store(config.TelemetryConfig{
			Mode: tc.Mode, Size: tc.Payload, Protocol: tc.Protocol, Sensors: tc.Sensors,
		})
		
		// 2. Sincronizza le goroutine (attiva i sensori necessari)
		syncSensors(tc.Sensors, clients, httpClient)

		// 3. Esecuzione Fasi
		log.Println("[1/3] Fase: Warm-up (30s)...")
		time.Sleep(30 * time.Second)

		log.Println("[2/3] Fase: Steady State (180s) - RACCOLTA DATI...")
		// NOTA: Se hai una funzione per resettare le metriche, chiamala qui
		time.Sleep(180 * time.Second)

		log.Println("[3/3] Fase: Cool-down (30s)...")
		syncSensors(0, clients, httpClient) // Spegne tutto per pulire i buffer
		time.Sleep(30 * time.Second)
		
		log.Printf(">>> TEST CASE COMPLETATO <<<\n")
	}

	log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	log.Println("!!! BENCHMARK SUITE TERMINATA !!!")
	log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
}

func updateGatewayConfig(client *http.Client, tc TestCase) {
    // Creiamo l'oggetto config da inviare
    cfg := config.TelemetryConfig{
        Mode:     tc.Mode,
        Size:     tc.Payload,
        Sensors:  tc.Sensors,
        Protocol: tc.Protocol,
    }

    // Trasformiamo in JSON
    jsonData, _ := json.Marshal(cfg)
    
    // Inviamo al Gateway (endpoint che abbiamo creato prima: /set-config)
    resp, err := client.Post(setConfigEndpoint, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        log.Printf("Errore notifica benchmark al gateway: %v", err)
        return
    }
    defer resp.Body.Close()
}