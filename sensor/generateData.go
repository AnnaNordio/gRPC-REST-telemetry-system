package main

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	// Usiamo l'alias 'pb' per il tuo pacchetto proto generato
	pb "telemetry-bench/proto"
)

func generateData(size string) *pb.SensorData {
	data := &pb.SensorData{
		SensorId:    "sensor_1",
		Temperature: 20.0 + rand.Float32()*10.0,
		Humidity:    40.0 + rand.Float32()*20.0,
		Timestamp:   time.Now().UnixMicro(),
	}

	switch size {
	case "medium":
		// strings.Repeat richiede il pacchetto "strings"
		data.PayloadContent = strings.Repeat("m", 10240) // 10KB
	case "large":
		data.PayloadContent = strings.Repeat("l", 102400) // 100KB
	case "nested":
		for i := 0; i < 50; i++ {
			// strconv richiede il pacchetto "strconv"
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
		// Caso "small" o default: payload minimo
		data.PayloadContent = "small_payload"
	}

	return data
}