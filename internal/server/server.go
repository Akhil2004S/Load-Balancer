package server

import (
	"loadBalancer/internal/health"
)

type ServerData struct {
	Id              int
	Address         string
	HealthState     health.HealthStatus
	TotalReqs       int
	SuccessReqs     int
	SuccessRate     float64
	ActiveConn      int
	AvgResponseTime float64
}

func CreateServer(id int, url string) *ServerData {
	return &ServerData{
		Id:      id,
		Address: url,
	}
}
