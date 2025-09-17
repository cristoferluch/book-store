package entity

type HealthResponse struct {
	Database string `json:"database,omitempty"`
	Version  string `json:"version"`
	Uptime   string `json:"uptime"`
}
