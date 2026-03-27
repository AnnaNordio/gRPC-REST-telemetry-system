package main

import (
    "io"
    "net/http"
)

func fetchConfig(client *http.Client) (string, string) {
    mode := fetchValue(client, modeEndpoint, "polling")
    size := fetchValue(client, sizeEndpoint, "small")
    return mode, size
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