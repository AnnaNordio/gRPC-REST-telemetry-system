package model

import "time"

// SensorReading è la nostra struttura dati interna "neutra"
type SensorReading struct {
	SensorID    string    `json:"sensor_id"`
	Temperature float32   `json:"temperature"`
	Humidity    float32   `json:"humidity"`
	Timestamp   time.Time `json:"timestamp"`
}

// Qui potresti aggiungere metodi di calcolo o validazione
func (s *SensorReading) IsValid() bool {
	return s.Temperature > -50 && s.Temperature < 100
}

