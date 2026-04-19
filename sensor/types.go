package main

const (
    gatewayBaseUrl  = "http://gateway:8080"
    gatewayRestAddr = gatewayBaseUrl + "/telemetry"
    gatewayGrpcAddr = "gateway:50051"
    configEndpoint    = gatewayBaseUrl + "/get-config"
    setConfigEndpoint = gatewayBaseUrl + "/set-config"
)