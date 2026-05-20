package server

type HealthStatus int

const (
	Unhealthy HealthStatus = iota
	Healthy
	Evaluating
)

var HealthState = map[HealthStatus]string{
	Healthy:    "Healthy",
	Unhealthy:  "Unhealthy",
	Evaluating: "Under Evaluation",
}
