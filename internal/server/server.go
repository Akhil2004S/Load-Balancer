package server

import (
	"loadBalancer/internal/health"
	"time"
)

type ServerData struct {
	Id              int
	Address         string
	HealthState     health.HealthStatus
	TotalReqs       int
	SuccessReqs     int
	ActiveConn      int
	AvgResponseTime time.Time
	Weight          float64
	UnhealthyCount  int
}

func CreateServer(id int, url string) *ServerData {
	return &ServerData{
		Id:      id,
		Address: url,
	}
}
