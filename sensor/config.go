package main

import (
    "io"
    "net/http"
)

func fetchConfig(client *http.Client) (string, string, string) {
    mode := fetchValue(client, modeEndpoint, "polling")
    size := fetchValue(client, sizeEndpoint, "small")
    sensors := fetchValue(client, sensorsEndpoint, "1")
    return mode, size, sensors
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