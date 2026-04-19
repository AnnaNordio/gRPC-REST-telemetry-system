package config

type TelemetryConfig struct {
    Mode     string `json:"mode"`
    Size     string `json:"size"`
    Sensors  int    `json:"sensors"`
    Protocol string `json:"protocol"`
}