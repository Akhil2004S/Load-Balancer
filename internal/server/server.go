package server

type ServerData struct {
	Id                 int
	Address            string
	HealthState        HealthStatus
	TotalReqs          int
	SuccessReqs        int
	SuccessRate        float64
	ActiveConn         int
	AvgResponseTime    float64
	HealthyThreshold   int
	UnhealthyThreshold int
	HealthyCount       int
	UnhealthyCount     int
}

func CreateServer(id int, url string) *ServerData {
	return &ServerData{
		Id:                 id,
		Address:            url,
		HealthyThreshold:   5,
		UnhealthyThreshold: 10,
	}
}
