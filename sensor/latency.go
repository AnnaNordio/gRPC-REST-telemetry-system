package main

var lastLatRest, lastLatGrpc float64

func updateLatency(protocol string, value float64) {
    if protocol == "REST" {
        lastLatRest = value
    } else {
        lastLatGrpc = value
    }
}