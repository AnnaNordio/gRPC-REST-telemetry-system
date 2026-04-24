package main

import (
	"math/rand"
	"strconv"
	pb "telemetry-bench/proto"
	"time"
)

func generateData(size string) *pb.SensorData {
	data := &pb.SensorData{
		SensorId:    "sensor_1",
		Temperature: 25.5,
		Humidity:    60.2,
		Timestamp:   time.Now().UnixMicro(),
	}

	switch size {
	case "medium":
		for i := 0; i < 200; i++ {
			data.Tags = append(data.Tags, "tag_category_name_"+strconv.Itoa(i))
		}

	case "large":
		data.ReadingsHistory = make(map[string]float32)
		for i := 0; i < 2000; i++ {
			data.ReadingsHistory["ts_"+strconv.Itoa(i)] = rand.Float32() * 100
		}

	case "nested":
		for i := 0; i < 100; i++ {
			data.Details = append(data.Details, &pb.NestedDetail{
				Key:   "parameter_" + strconv.Itoa(i),
				Value: "value_hash_7b9c1d2e3f4a5",
				Metadata: map[string]string{
					"node_id": "cluster_alpha_north",
					"status":  "active",
					"version": "1.2.4-stable",
				},
			})
		}

	default:
		data.Tags = []string{"stable", "indoor"}
	}

	return data
}
