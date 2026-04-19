package main

import (
	"io"
	"net/http"
	"strconv"
)

// RemoteConfig definisce i parametri operativi del benchmark
type RemoteConfig struct {
	Mode     string
	Size     string
	Sensors  int
	Protocol string
}

// fetchFullConfig interroga gli endpoint e restituisce la nuova configurazione.
// Se un valore non è disponibile, usa quello attuale per evitare reset improvvisi.
func fetchFullConfig(client *http.Client) RemoteConfig {
	// Recuperiamo lo stato attuale per i fallback
	current := globalConfig.Load().(RemoteConfig)

	m := fetchValue(client, modeEndpoint, current.Mode)
	sz := fetchValue(client, sizeEndpoint, current.Size)
	p := fetchValue(client, protocolEndpoint, current.Protocol)
	
	// Gestione sicura per il numero di sensori
	sStr := fetchValue(client, sensorsEndpoint, strconv.Itoa(current.Sensors))
	sInt, err := strconv.Atoi(sStr)
	if err != nil {
		sInt = current.Sensors
	}

	return RemoteConfig{
		Mode:     m,
		Size:     sz,
		Sensors:  sInt,
		Protocol: p,
	}
}

func fetchValue(client *http.Client, url string, defaultValue string) string {
	resp, err := client.Get(url)
	if err != nil {
		return defaultValue
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil || len(body) == 0 {
		return defaultValue
	}
	return string(body)
}