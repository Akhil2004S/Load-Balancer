package health

type HealthStatus int

const (
	healthy HealthStatus = iota
	unhealthy
	evaluating
)

var healthState = map[HealthStatus]string{
	healthy:    "Healthy",
	unhealthy:  "Unhealthy",
	evaluating: "Under Evaluation",
}
