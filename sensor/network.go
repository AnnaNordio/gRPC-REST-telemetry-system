package main

import (
    "bytes"
    "context"
    "encoding/json"
    "io"
    "net/http"
    "time"
	"google.golang.org/protobuf/proto"
	"fmt"

    pb "telemetry-bench/proto"
)

func sendRest(client *http.Client, data *pb.SensorData) {
	// 1. Serializziamo in JSON
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Errore marshal JSON: %v", err)
		return
	}

	// 2. Creiamo la richiesta
	req, _ := http.NewRequest("POST", gatewayRestAddr, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	// 3. Invio
	resp, err := client.Do(req)
	if err == nil {
		// Leggiamo la risposta per liberare la connessione (Keep-Alive)
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	} else {
		fmt.Printf("Errore invio REST: %v", err)
	}
}

func executeStreaming(client *http.Client, stream pb.TelemetryService_StreamDataClient, data *pb.SensorData) {
	data.Timestamp = time.Now().UnixMicro()

	// --- LOGICA PESO gRPC STREAM ---
	// Calcoliamo quanto pesa il messaggio protobuf
	pBytes, _ := proto.Marshal(data)
	data.PayloadBytes = int64(len(pBytes))
	
	// Invio allo stream
	if err := stream.Send(data); err != nil {
		fmt.Printf("Errore Stream Send: %v", err)
	}

	// Invio REST asincrono (per confronto costante)
	// Passiamo una copia per evitare che il calcolo di PayloadBytes gRPC 
	// sovrascriva quello che serve a REST
	dataCopy := *data 
	go sendRest(client, &dataCopy)
}

func executePolling(client *http.Client, grpcClient pb.TelemetryServiceClient, data *pb.SensorData) {
	data.Timestamp = time.Now().UnixMicro()

	// 1. Esecuzione REST (il peso viene calcolato dentro sendRest)
	dataRest := *data
	go sendRest(client, &dataRest)

	// 2. Esecuzione gRPC Unary
	go func(d pb.SensorData) {
		// --- LOGICA PESO gRPC UNARY ---
		// Calcoliamo il peso protobuf specifico per questa chiamata
		pBytes, _ := proto.Marshal(&d)
		d.PayloadBytes = int64(len(pBytes))

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err := grpcClient.SendData(ctx, &d)
		if err != nil {
			fmt.Printf("Errore gRPC Unary: %v", err)
		}
	}(*data)
}