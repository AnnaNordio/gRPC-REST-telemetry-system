package main

type DashboardPayload struct {
    AvgRestSize float64 `json:"avg_rest_size"`
    AvgGrpcSize float64 `json:"avg_grpc_size"`
}

var (
    sumSizeRest, countSizeRest float64
    sumSizeGrpc, countSizeGrpc float64
)

func savePayload(protocol string, size int) {
    metricsMu.Lock()
    defer metricsMu.Unlock()

    if protocol == "REST" {
        sumSizeRest += float64(size)
        countSizeRest++
    } else {
        sumSizeGrpc += float64(size)
        countSizeGrpc++
    }
}