package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"github.com/anna/iot-dual-stack/gen/sensor";
)

func main() {
	// Setup gRPC Client
	conn, _ := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	grpcClient := sensor.NewTelemetryServiceClient(conn)

	// Dati finti
	data := &sensor.TelemetryData{
		SensorId:    "sensor-01",
		Temperature: 24.5,
		Humidity:    60.2,
		Timestamp:   time.Now().Unix(),
	}

	for {
		// 1. Test gRPC
		start := time.Now()
		_, err := grpcClient.SendData(context.Background(), data)
		if err == nil {
			log.Printf("gRPC Call riuscita in %v", time.Since(start))
		}

		// 2. Test REST
		start = time.Now()
		jsonData, _ := json.Marshal(data)
		resp, err := http.Post("http://localhost:8080/telemetry", "application/json", bytes.NewBuffer(jsonData))
		if err == nil {
			log.Printf("REST Call riuscita in %v (Size: %d)", time.Since(start), len(jsonData))
			resp.Body.Close()
		}

		time.Sleep(2 * time.Second)
	}
}

