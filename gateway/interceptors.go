package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	pb "telemetry-bench/proto"
)

// --- gRPC UNARY INTERCEPTOR ---
func grpcSizeInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 1. Filtriamo solo il metodo interessato
	if info.FullMethod != "/telemetry.TelemetryService/SendData" {
		return handler(ctx, req)
	}

	// 2. Verifichiamo che il dato sia quello atteso
	data, ok := req.(*pb.SensorData)
	if !ok {
		return handler(ctx, req)
	}

	// 3. Calcolo Overhead Headers (Metadata HTTP/2)
	md, _ := metadata.FromIncomingContext(ctx)
	var hSize int64
	for k, v := range md {
		hSize += int64(len(k))
		for _, s := range v {
			hSize += int64(len(s))
		}
	}

	// 4. Esecuzione dell'handler
	resp, err := handler(ctx, req)

	// 5. CHIAMATA UNICA: Salva Latenza + Payload + Overhead
	// Usiamo data.Timestamp per calcolare la latenza reale
	saveAllMetrics("gRPC", data.Timestamp, data.PayloadBytes, hSize)

	log.Printf("[gRPC Unary] Processed: Payload %d, Overhead %d", data.PayloadBytes, hSize)
	return resp, err
}

// --- REST MIDDLEWARE ---
func restSizeMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Calcolo Overhead Headers (HTTP/1.1)
		var hSize int64
		for k, v := range r.Header {
			hSize += int64(len(k))
			for _, s := range v {
				hSize += int64(len(s))
			}
		}

		// 2. Lettura del Body per estrarre Timestamp e dimensione Payload
		body, err := io.ReadAll(r.Body)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		pSize := int64(len(body))

		// Estraiamo il timestamp dal JSON per la latenza
		var sensorData struct {
			Timestamp int64 `json:"timestamp"`
		}
		json.Unmarshal(body, &sensorData)

		// Fondamentale: ripristiniamo il body per l'handler handleTelemetry
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		// 3. CHIAMATA UNICA: Salva Latenza + Payload + Overhead
		saveAllMetrics("REST", sensorData.Timestamp, pSize, hSize)

		log.Printf("[REST] Processed: Payload %d, Overhead %d", pSize, hSize)
		
		// 4. Continua l'esecuzione (handleTelemetry NON deve più chiamare saveMetric)
		next.ServeHTTP(w, r)
	}
}