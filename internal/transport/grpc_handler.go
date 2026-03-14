package transport

import (
	"context"
	"log"
	"time"

	"github.com/anna/iot-dual-stack/gen/sensor"
)

// GRPCHandler gestisce le chiamate gRPC in arrivo
type GRPCHandler struct {
	sensor.UnimplementedTelemetryServiceServer
}

func (h *GRPCHandler) SendData(ctx context.Context, req *sensor.TelemetryData) (*sensor.Reply, error) {
	// 1. Misuriamo l'arrivo (per la tesi)
	start := time.Now()
	
	log.Printf("[gRPC-Transport] Ricevuto da %s: Temp %.2f", req.SensorId, req.Temperature)

	// In un caso reale, qui chiameresti la logica di business o inoltreresti al backend
	
	elapsed := time.Since(start)
	return &sensor.Reply{
		Status:     "OK",
		ServerTime: time.Now().Unix(),
	}, nil
}

