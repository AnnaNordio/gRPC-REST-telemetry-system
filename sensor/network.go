package main

import (
    "bytes"
    "context"
    "encoding/json"
    "io"
    "net/http"
    "time"
    "log"

    pb "telemetry-bench/proto"
)

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

func executeStreaming(client *http.Client, stream pb.TelemetryService_StreamDataClient, data *pb.SensorData) {
    data.Timestamp = time.Now().UnixMicro()
    
    // Invio gRPC Stream
    _ = stream.Send(data)
    
    // Invio REST asincrono
    go sendRest(client, data)
}

// executePolling modificata per essere più robusta
func executePolling(client *http.Client, grpcClient pb.TelemetryServiceClient, data *pb.SensorData) {
	data.Timestamp = time.Now().UnixMicro()
	
	// REST asincrono
	go sendRest(client, data)
	
	// gRPC Unary
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err := grpcClient.SendData(ctx, data)
	if err != nil {
		log.Printf("Errore gRPC Polling: %v", err)
	}
}