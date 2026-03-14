package transport

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/anna/iot-dual-stack/gen/sensor"
)

// RestHandler gestisce le chiamate REST (JSON)
func RestHandler(w http.ResponseWriter, r *http.Request) {
	// Misuriamo il tempo di gestione
	start := time.Now()

	if r.Method != http.MethodPost {
		http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
		return
	}

	// Leggiamo il body (utile per misurare la dimensione del payload)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Errore lettura body", http.StatusBadRequest)
		return
	}

	// Decodifica JSON nella struct generata dal Proto (per coerenza)
	var data sensor.TelemetryData
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("[REST-Transport] Errore unmarshal: %v", err)
		http.Error(w, "JSON non valido", http.StatusBadRequest)
		return
	}

	log.Printf("[REST-Transport] Ricevuto da %s: %d bytes in %v", 
		data.SensorId, len(body), time.Since(start))

	// Risposta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "OK",
		"received_size": fmt.Sprintf("%d bytes", len(body)),
	})
}

