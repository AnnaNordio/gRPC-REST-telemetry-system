package main

import (
    "io"
    "net/http"
    "strconv"
)

type RemoteConfig struct {
    Mode    string
    Size    string
    Sensors int
    Protocol string
}

func fetchFullConfig(client *http.Client) RemoteConfig {
    m := fetchValue(client, modeEndpoint, "polling")
    sz := fetchValue(client, sizeEndpoint, "small")
    sStr := fetchValue(client, sensorsEndpoint, "1")
    p := fetchValue(client, protocolEndpoint, "both")
    
    sInt, _ := strconv.Atoi(sStr) 
    
    return RemoteConfig{
        Mode:    m,
        Size:    sz,
        Sensors: sInt,
        Protocol: p,
    }
}

func fetchValue(client *http.Client, url string, defaultValue string) string {
    resp, err := client.Get(url)
    if err != nil {
        return defaultValue
    }
    defer resp.Body.Close()
    body, _ := io.ReadAll(resp.Body)
    return string(body)
}