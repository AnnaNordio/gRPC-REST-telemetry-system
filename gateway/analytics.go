package main

import (
    "encoding/json"
    "time"
    "sort"
    "net/http"
    "google.golang.org/grpc/metadata"  
    "google.golang.org/protobuf/proto"
)

const aggregateWindow = 200

func getDashboardData() DashboardResponse {
    metricsMu.Lock()
    defer metricsMu.Unlock()

    // 1. Prendiamo gli ultimi 200 campioni per protocollo (Indipendentemente dal sensore)
    restLats := getLastNLats(history, "REST", aggregateWindow)
    grpcLats := getLastNLats(history, "gRPC", aggregateWindow)

    // 2. Calcolo Latenza Media Aggregata (su finestra mobile)
    avgRest := calculateAverage(restLats)
    avgGrpc := calculateAverage(grpcLats)

    // 3. Calcolo P99 Aggregato
    p99Rest := calculatePercentile(restLats, 0.99)
    p99Grpc := calculatePercentile(grpcLats, 0.99)

    return DashboardResponse{
        History:           history, // La history intera serve per disegnare le linee del grafico
        AvgRest:           avgRest, // I "numeroni" nelle card
        AvgGrpc:           avgGrpc,
        P99Rest:           p99Rest,
        P99Grpc:           p99Grpc,
        TotalPayloadRest:  totalPayloadRest,
        TotalOverheadRest: totalOverheadRest,
        TotalPayloadGrpc:  totalPayloadGrpc,
        TotalOverheadGrpc: totalOverheadGrpc,
        LastGrpcTSRaw:     lastGlobalGrpcTS,
    }
}

// Helper per isolare le latenze dell'ultimo periodo
func getLastNLats(h []Metric, protocol string, n int) []float64 {
    var lats []float64
    for i := len(h) - 1; i >= 0 && len(lats) < n; i-- {
        if h[i].Protocol == protocol {
            lats = append(lats, h[i].LatencyMs)
        }
    }
    return lats
}

func calculateAverage(lats []float64) float64 {
    if len(lats) == 0 { return 0 }
    var sum float64
    for _, l := range lats { sum += l }
    return sum / float64(len(lats))
}

func calculateLatency(sensorTS int64) float64 {
    now := time.Now().UnixMicro()
    return float64(now - sensorTS)
}

func safeAvg(sum float64, count int) float64 {
    if count == 0 {
        return 0
    }
    return sum / float64(count)
}

// Calcolo matematico del percentile (P99)
func calculatePercentile(latencies []float64, percentile float64) float64 {
    if len(latencies) == 0 {
        return 0
    }
    
    // Copia e ordina per non alterare la history originale
    sorted := make([]float64, len(latencies))
    copy(sorted, latencies)
    sort.Float64s(sorted)

    index := int(float64(len(sorted)-1) * percentile)
    return sorted[index]
}

// getJsonSize restituisce la dimensione in byte del payload serializzato in JSON
func getJsonSize(v interface{}) int64 {
	b, err := json.Marshal(v)
	if err != nil {
		return 0
	}
	return int64(len(b))
}

// getProtoSize restituisce la dimensione in byte del payload serializzato in Protobuf
func getProtoSize(m proto.Message) int64 {
	b, err := proto.Marshal(m)
	if err != nil {
		return 0
	}
	return int64(len(b))
}

// Supporto per calcolare i byte reali degli header REST
func calculateHTTPOverhead(req *http.Request) int64 {
    var size int64
    
    // 1. Request Line reale: "GET /api/data HTTP/1.1" + \r\n
    // Calcola esattamente la lunghezza delle stringhe reali della richiesta
    size += int64(len(req.Method) + 1 + len(req.URL.RequestURI()) + 1 + len(req.Proto) + 2)
    
    // 2. Headers reali
    for name, values := range req.Header {
        for _, v := range values {
            // "Name: value\r\n"
            size += int64(len(name) + 2 + len(v) + 2)
        }
    }
    
    // 3. L'ultima riga vuota che separa headers dal body (\r\n)
    size += 2
    
    return size
}

// Supporto per calcolare i byte reali dei metadati gRPC
func calculateGRPCOverhead(md metadata.MD) int64 {
	var size int64
	for k, vs := range md {
		size += int64(len(k))
		for _, v := range vs {
			size += int64(len(v))
		}
	}
	return size
}