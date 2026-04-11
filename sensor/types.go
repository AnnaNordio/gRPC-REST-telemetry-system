package main

const (
    gatewayBaseUrl  = "http://gateway:8080"
    gatewayRestAddr = gatewayBaseUrl + "/telemetry"
    gatewayGrpcAddr = "gateway:50051"
    modeEndpoint    = gatewayBaseUrl + "/get-mode"
    sizeEndpoint    = gatewayBaseUrl + "/get-size"
    sensorsEndpoint = gatewayBaseUrl + "/get-sensors"
    protocolEndpoint = gatewayBaseUrl + "/get-protocol"
)