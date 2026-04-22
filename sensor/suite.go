package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	pb "telemetry-bench/proto"
	"telemetry-bench/pkg/config"
)

type TestCase struct {
	Sensors  int
	Mode     string
	Payload  string
	Protocol string
}

// ----------------------------
// BUILD GROUPED SUITE
// ----------------------------
func buildGroupedSuite() [][]TestCase {
	sensors := []int{1, 10, 50, 100}
	modes := []string{"polling", "streaming"}
	payloads := []string{"small", "medium", "large", "nested"}
	protocols := []string{"grpc", "rest"}

	groups := make([][]TestCase, len(sensors))

	for i, s := range sensors {
		var group []TestCase

		for _, m := range modes {
			for _, p := range payloads {
				for _, proto := range protocols {
					group = append(group, TestCase{
						Sensors:  s,
						Mode:     m,
						Payload:  p,
						Protocol: proto,
					})
				}
			}
		}

		groups[i] = group
	}

	return groups
}

// ----------------------------
// SHUFFLE EACH GROUP
// ----------------------------
func shuffleGroups(groups [][]TestCase) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range groups {
		r.Shuffle(len(groups[i]), func(a, b int) {
			groups[i][a], groups[i][b] = groups[i][b], groups[i][a]
		})
	}
}

// ----------------------------
// INTERLEAVE GROUPS
// ----------------------------
func interleave(groups [][]TestCase) []TestCase {
	var result []TestCase

	maxLen := 0
	for _, g := range groups {
		if len(g) > maxLen {
			maxLen = len(g)
		}
	}

	for i := 0; i < maxLen; i++ {
		for _, g := range groups {
			if i < len(g) {
				result = append(result, g[i])
			}
		}
	}

	return result
}

// ----------------------------
// MAIN BENCHMARK RUNNER
// ----------------------------
func runBenchmarkSuite(clients []pb.TelemetryServiceClient, httpClient *http.Client) {

	groups := buildGroupedSuite()

	// 🔀 shuffle interno per evitare pattern fissi
	shuffleGroups(groups)

	// 🔁 interleave tra livelli di carico
	suite := interleave(groups)

	for i, tc := range suite {

		log.Printf("\n>>> [%d/%d] TEST CASE: Sensors:%d | Mode:%s | Size:%s | Protocol:%s <<<",
			i+1, len(suite),
			tc.Sensors, tc.Mode, tc.Payload, tc.Protocol,
		)

		updateGatewayConfig(httpClient, tc)

		globalConfig.Store(config.TelemetryConfig{
			Mode:     tc.Mode,
			Size:     tc.Payload,
			Protocol: tc.Protocol,
			Sensors:  tc.Sensors,
		})

		syncSensors(tc.Sensors, clients, httpClient)

		log.Println("[1/3] Warm-up (30s)")
		time.Sleep(30 * time.Second)

		log.Println("[2/3] Steady state (180s)")
		time.Sleep(180 * time.Second)

		log.Println("[3/3] Cool-down (30s)")
		syncSensors(0, clients, httpClient)
		time.Sleep(30 * time.Second)

		log.Println(">>> TEST COMPLETATO <<<")
	}

	log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	log.Println("!!! BENCHMARK SUITE TERMINATA !!!")
	log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
}

// ----------------------------
// UPDATE GATEWAY CONFIG
// ----------------------------
func updateGatewayConfig(client *http.Client, tc TestCase) {

	cfg := config.TelemetryConfig{
		Mode:     tc.Mode,
		Size:     tc.Payload,
		Sensors:  tc.Sensors,
		Protocol: tc.Protocol,
	}

	jsonData, _ := json.Marshal(cfg)

	resp, err := client.Post(
		setConfigEndpoint,
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		log.Printf("Errore notifica benchmark al gateway: %v", err)
		return
	}
	defer resp.Body.Close()
}