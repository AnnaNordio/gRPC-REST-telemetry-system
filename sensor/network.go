package main

import (
    "bytes"
    "context"
    "encoding/json"
    "io"
    "net/http"
    "time"

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

func executePolling(client *http.Client, grpcClient pb.TelemetryServiceClient, data *pb.SensorData) {
    data.Timestamp = time.Now().UnixMicro()

    // Entrambi asincroni per massima velocità del loop
    go sendRest(client, data)
    go func() {
        _, _ = grpcClient.SendData(context.Background(), data)
    }()
}