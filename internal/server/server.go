package server

import (
	"loadBalancer/internal/health"
	"time"
)

type ServerData struct {
	Address         string
	HealthState     health.HealthStatus
	TotalReqs       int
	SuccessReqs     int
	ActiveConn      int
	AvgResponseTime time.Time
	Weight          float64
	UnhealthyCount  int
}

func CreateServer(url string) *ServerData {
	return &ServerData{
		Address: url,
	}
}
