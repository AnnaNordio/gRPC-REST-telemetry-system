package metrics

import (
	"fmt"
	"time"
)

// Metric raccoglie i dati di un singolo invio
type Metric struct {
	Protocol    string  // "REST" o "gRPC"
	Bytes       int     // Dimensione del pacchetto
	Duration    int64   // Durata in microsecondi
	Success     bool
}

// GetCSVHeader restituisce l'intestazione per il tuo file Excel
func GetCSVHeader() string {
	return "Protocol,Bytes,Duration_us,Success\n"
}

// ToCSV trasforma la metrica in una riga CSV
func (m Metric) ToCSV() string {
	return fmt.Sprintf("%s,%d,%d,%v\n", m.Protocol, m.Bytes, m.Duration, m.Success)
}

// CalculateSavings confronta le dimensioni e calcola il risparmio percentuale
func CalculateSavings(restSize, grpcSize int) float64 {
	if restSize == 0 {
		return 0
	}
	return float64(restSize-grpcSize) / float64(restSize) * 100
}