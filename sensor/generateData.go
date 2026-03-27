package main

import (
    "math/rand"
    "time"
    pb "telemetry-bench/proto"
)

func generateData() *pb.SensorData {
    return &pb.SensorData{
        SensorId:    "sensor-01",
        Temperature: 20 + rand.Float32()*10,
        Timestamp:   string(time.Now().UnixMilli()),
        LatencyRest: lastLatRest,
        LatencyGrpc: lastLatGrpc,
    }
}