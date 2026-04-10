package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"	
	pb "telemetry-bench/proto"
)

// Middleware per gestire le richieste Cross-Origin
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Restituisce i dati per la dashboard (filtrando la history per REST)
func handleResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fullData := getDashboardData()

	// Filtriamo la history per mostrare solo i dati REST nel grafico dedicato
	var restOnlyHistory []Metric
	for _, m := range fullData.History {
		if m.Protocol == "REST" {
			restOnlyHistory = append(restOnlyHistory, m)
		}
	}

	response := struct {
		History []Metric `json:"history"`
		AvgRest float64  `json:"avg_rest"`
		P99Rest float64  `json:"p99_rest"`
		PayloadSize int64 `json:"payload_size"`
		Overhead int64 `json:"overhead_size"`
		Throughput float64 `json:"throughput_rest"`
		MarshalTime float64 `json:"marshal_time_rest"`
	}{
		History: restOnlyHistory,
		AvgRest: fullData.AvgRest,
		P99Rest: fullData.P99Rest,
		PayloadSize: fullData.TotalPayloadRest,
        Overhead: fullData.TotalOverheadRest,
		Throughput: fullData.ThroughputRest,
		MarshalTime: fullData.MarshalAvgRest,
	}

	json.NewEncoder(w).Encode(response)
}

// Riceve i dati dai sensori via REST
func handleTelemetry(w http.ResponseWriter, r *http.Request) {
	var data pb.SensorData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	SaveRestMetrics(&data, r)
	w.WriteHeader(http.StatusOK)
}

// Cambia la modalità (polling vs streaming)
func handleSetMode(w http.ResponseWriter, r *http.Request) {
	newMode := r.URL.Query().Get("mode")
	if newMode != "" && (newMode == "polling" || newMode == "streaming") {
		if newMode != currentMode {
			currentMode = newMode
			resetStats() // Resetta le metriche quando cambia il paradigma
		}
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Modalità non valida", http.StatusBadRequest)
	}
}

// Restituisce la modalità attuale
func handleGetMode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, currentMode)
}

// Cambia la dimensione del payload simulato
func handleSetSize(w http.ResponseWriter, r *http.Request) {
	newSize := r.URL.Query().Get("size")
	if newSize != "" {
		if newSize != currentSize {
			currentSize = newSize
			resetStats() 
		}
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Dimensione non valida", http.StatusBadRequest)
	}
}

// Restituisce la dimensione del payload attuale
func handleGetSize(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, currentSize)
}

func handleSetSensors(w http.ResponseWriter, r *http.Request) {
	newCount := r.URL.Query().Get("count")
	if newCount != "" {
		if newCount != currentSensors {
			currentSensors = newCount
			resetStats()
		}
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Numero sensori non valido", http.StatusBadRequest)
	}
}

// Restituisce il numero di sensori attuale
func handleGetSensors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, currentSensors)
}

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true }, 
}

func handleWS(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return
    }
    defer conn.Close()

    for {
        if _, _, err := conn.ReadMessage(); err != nil {
            break 
        }
    }
}

func handleReset(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Metodo non consentito", http.StatusMethodNotAllowed)
        return
    }
    
    resetStats() 
    w.WriteHeader(http.StatusOK)
}