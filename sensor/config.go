package main

import (
	"encoding/json"
	"log"
	"net/http"
	"telemetry-bench/pkg/config"
)

func fetchFullConfig(client *http.Client) config.TelemetryConfig {
	// fallback
	current := globalConfig.Load().(config.TelemetryConfig)

	resp, err := client.Get(configEndpoint)
	if err != nil {
		log.Printf("Errore fetch config: %v", err)
		return current
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return current
	}

	var newConfig config.TelemetryConfig
	if err := json.NewDecoder(resp.Body).Decode(&newConfig); err != nil {
		log.Printf("Errore decode config: %v", err)
		return current
	}

	return newConfig
}
