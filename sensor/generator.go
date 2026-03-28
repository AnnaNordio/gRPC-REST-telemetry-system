package main

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
	pb "telemetry-bench/proto"
)

// generateData crea un oggetto SensorData con dimensioni variabili basate sul parametro size
func generateData(size string) *pb.SensorData {
	data := &pb.SensorData{
		SensorId:    "sensor_1",
		Temperature: 20.0 + rand.Float32()*10.0,
		Humidity:    40.0 + rand.Float32()*20.0,
		Timestamp:   time.Now().UnixMicro(),
	}

	switch size {
	case "medium":
		data.PayloadContent = strings.Repeat("m", 10240) // 10KB
	case "large":
		data.PayloadContent = strings.Repeat("l", 102400) // 100KB
	case "nested":
		// Genera dati complessi per testare la serializzazione di strutture annidate
		for i := 0; i < 50; i++ {
			data.Details = append(data.Details, &pb.NestedDetail{
				Key:   "attr_" + strconv.Itoa(i),
				Value: "value_data_point",
				Metadata: map[string]string{
					"unit":   "celsius",
					"status": "ok",
				},
			})
		}
	default:
		// Caso "small": payload minimo
		data.PayloadContent = "small_payload"
	}

	return data
}