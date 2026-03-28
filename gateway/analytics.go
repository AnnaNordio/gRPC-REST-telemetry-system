package main

import (
    "encoding/json"
    "time"
    "sort"
    "net/http"
    "google.golang.org/grpc/metadata"  
    "google.golang.org/protobuf/proto"
)

func getDashboardData() DashboardResponse {
    metricsMu.Lock()
    defer metricsMu.Unlock()

    var restLats, grpcLats []float64
    var sumR, sumG float64

    for _, m := range history {
        if m.Protocol == "REST" {
            restLats = append(restLats, m.LatencyMs)
            sumR += m.LatencyMs
        } else {
            grpcLats = append(grpcLats, m.LatencyMs)
            sumG += m.LatencyMs
        }
    }

    return DashboardResponse{
        History:       history,
        TotalPayloadRest:  totalPayloadRest,
        TotalOverheadRest: totalOverheadRest,
        TotalPayloadGrpc:  totalPayloadGrpc,
        TotalOverheadGrpc: totalOverheadGrpc,
        AvgRest:       safeAvg(sumR, len(restLats)),
        AvgGrpc:       safeAvg(sumG, len(grpcLats)),
        P99Rest:       calculatePercentile(restLats, 0.99),
        P99Grpc:       calculatePercentile(grpcLats, 0.99),
        LastGrpcTSRaw: lastGlobalGrpcTS,
    }
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