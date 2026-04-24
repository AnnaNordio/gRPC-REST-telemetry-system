package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"telemetry-bench/pkg/config"
	pb "telemetry-bench/proto"

	"github.com/gorilla/websocket"
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

	var restOnlyHistory []Metric
	for _, m := range fullData.History {
		if m.Protocol == "REST" {
			restOnlyHistory = append(restOnlyHistory, m)
		}
	}

	response := struct {
		History     []Metric `json:"history"`
		AvgRest     float64  `json:"avg_rest"`
		P99Rest     float64  `json:"p99_rest"`
		PayloadSize int64    `json:"payload_size"`
		Overhead    int64    `json:"overhead_size"`
		Throughput  float64  `json:"throughput_rest"`
		MarshalTime float64  `json:"marshal_time_rest"`
	}{
		History:     restOnlyHistory,
		AvgRest:     fullData.AvgRest,
		P99Rest:     fullData.P99Rest,
		PayloadSize: fullData.TotalPayloadRest,
		Overhead:    fullData.TotalOverheadRest,
		Throughput:  fullData.ThroughputRest,
		MarshalTime: fullData.MarshalAvgRest,
	}

	json.NewEncoder(w).Encode(response)
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
	if newMode == "polling" || newMode == "streaming" {
		metricsMu.Lock()
		if activeConfig.Mode != newMode {
			activeConfig.Mode = newMode
			metricsMu.Unlock()
			resetStats()
		} else {
			metricsMu.Unlock()
		}
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Modalità non valida", http.StatusBadRequest)
	}
}

// Cambia la dimensione del payload
func handleSetSize(w http.ResponseWriter, r *http.Request) {
	newSize := r.URL.Query().Get("size")
	if newSize != "" {
		metricsMu.Lock()
		if activeConfig.Size != newSize {
			activeConfig.Size = newSize
			metricsMu.Unlock()
			resetStats()
		} else {
			metricsMu.Unlock()
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handleSetSensors(w http.ResponseWriter, r *http.Request) {
	newCountStr := r.URL.Query().Get("count")
	newCount, err := strconv.Atoi(newCountStr)
	if err == nil {
		metricsMu.Lock()
		if activeConfig.Sensors != newCount {
			activeConfig.Sensors = newCount
			metricsMu.Unlock()
			resetStats()
		} else {
			metricsMu.Unlock()
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handleSetProtocol(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query().Get("p")
	if p != "" {
		metricsMu.Lock()
		if activeConfig.Protocol != p {
			activeConfig.Protocol = p
			metricsMu.Unlock()
			resetStats()
		} else {
			metricsMu.Unlock()
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handleGetMode(w http.ResponseWriter, r *http.Request) {
	metricsMu.Lock()
	mode := activeConfig.Mode
	metricsMu.Unlock()
	fmt.Fprint(w, mode)
}

func handleGetSensors(w http.ResponseWriter, r *http.Request) {
	metricsMu.Lock()
	s := activeConfig.Sensors
	metricsMu.Unlock()
	fmt.Fprint(w, s)
}

// Restituisce la dimensione del payload attuale
func handleGetSize(w http.ResponseWriter, r *http.Request) {
	metricsMu.Lock()
	size := activeConfig.Size
	metricsMu.Unlock()
	fmt.Fprint(w, size)
}

func handleGetProtocol(w http.ResponseWriter, r *http.Request) {
	metricsMu.Lock()
	protocol := activeConfig.Protocol
	metricsMu.Unlock()
	fmt.Fprint(w, protocol)
}

func handleSetConfig(w http.ResponseWriter, r *http.Request) {
	var newCfg config.TelemetryConfig
	if err := json.NewDecoder(r.Body).Decode(&newCfg); err != nil {
		http.Error(w, "JSON non valido", http.StatusBadRequest)
		return
	}

	metricsMu.Lock()
	activeConfig = newCfg
	metricsMu.Unlock()

	resetStats()
	w.WriteHeader(http.StatusOK)
}

func handleGetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	metricsMu.Lock()
	configToSend := activeConfig
	metricsMu.Unlock()

	err := json.NewEncoder(w).Encode(configToSend)
	if err != nil {
		http.Error(w, "Errore durante la codifica del JSON", http.StatusInternalServerError)
		return
	}
}
