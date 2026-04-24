package main

import (
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	pb "telemetry-bench/proto"
	"time"

	"google.golang.org/grpc/metadata"
)

// Canale per trasportare le metriche dai server al worker
// Capacità 10.000 per gestire picchi di 100 sensori a 10Hz
var metricsChan = make(chan Metric, 10000)
var historyBuffer = NewRingBuffer(1000)

func metricsWorker() {
	isBenchMode := os.Getenv("BENCH_MODE") == "true"
	if isBenchMode {
		os.MkdirAll("results", 0755)
	}

	writer := &MetricsWriter{}

	flushTicker := time.NewTicker(1 * time.Second)

	for {
		select {
		case m, ok := <-metricsChan:
			if !ok {
				return
			}

			// --- 1. Lettura Stato ---
			metricsMu.Lock()
			isWarmup := time.Now().Before(warmupUntil)
			cfg := activeConfig
			metricsMu.Unlock()

			if isWarmup {
				continue
			}

			// --- 2. Aggiornamento Statistiche  ---
			updateStats(m)

			// --- 3. Scrittura su File ---
			if isBenchMode {
				sSensors := strconv.Itoa(cfg.Sensors)
				writer.Write(m, cfg.Mode, cfg.Size, cfg.Protocol, sSensors)
			}

		case <-flushTicker.C:
			if writer.csvWriter != nil {
				writer.csvWriter.Flush()
			}
		}

	}
}

func updateStats(m Metric) {
	metricsMu.Lock()
	defer metricsMu.Unlock()

	if m.Protocol == "gRPC" {
		atomic.AddUint64(&msgCountGrpc, 1)
		totalPayloadGrpc += m.PayloadByte
		totalOverheadGrpc += m.OverheadByte
	} else {
		atomic.AddUint64(&msgCountRest, 1)
		totalPayloadRest += m.PayloadByte
		totalOverheadRest += m.OverheadByte
	}

	historyBuffer.Add(m)
}

func SaveRestMetrics(data *pb.SensorData, r *http.Request) {

	pSize, mTime := getJsonMetrics(data)
	hSize := calculateHTTPOverhead(r)
	lat := calculateLatency(data.Timestamp)

	metricsChan <- Metric{
		Protocol:     "REST",
		LatencyMs:    lat,
		PayloadByte:  pSize,
		OverheadByte: hSize,
		MarshalTime:  mTime,
		Timestamp:    time.Now().Format("15:04:05.000"),
	}
}

func SaveGrpcMetrics(data *pb.SensorData, md metadata.MD) {

	pSize, mTime := getProtoMetrics(data)
	hSize := 5 + calculateGRPCOverhead(md)
	lat := calculateLatency(data.Timestamp)

	metricsChan <- Metric{
		Protocol:     "gRPC",
		LatencyMs:    lat,
		PayloadByte:  pSize,
		OverheadByte: hSize,
		MarshalTime:  mTime,
		Timestamp:    time.Now().Format("15:04:05.000"),
	}
}

func resetStats() {
	for {
		select {
		case <-metricsChan:
		default:
			goto channelEmptied
		}
	}
channelEmptied:

	metricsMu.Lock()
	defer metricsMu.Unlock()

	warmupUntil = time.Now().Add(warmupDuration)

	historyBuffer.Reset()

	totalPayloadRest = 0
	totalOverheadRest = 0
	totalPayloadGrpc = 0
	totalOverheadGrpc = 0

	atomic.StoreUint64(&msgCountRest, 0)
	atomic.StoreUint64(&msgCountGrpc, 0)

	throughputRest = 0
	throughputGrpc = 0
}
