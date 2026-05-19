package algorithms

import (
	"fmt"
	"loadBalancer/internal/server"
	"math/rand/v2"
	"slices"
)

type LeastResponseTime struct {
	TotalRequests     int
	TotalResponseTime float64
}

func (lrt LeastResponseTime) NextServer(servers []*server.ServerData) *server.ServerData {
	if lrt.TotalRequests <= 10 {
		numServers := len(servers)
		randServer := rand.IntN(numServers)
		return servers[randServer]
	}
	var weights []float64
	for _, server := range servers {
		if lrt.TotalRequests <= 100 {
			weights = append(weights, calcWeight(1, lrt.TotalResponseTime, server))
		} else {
			alpha := float64(server.TotalReqs) / float64(lrt.TotalRequests)
			weights = append(weights, calcWeight(float64(alpha), lrt.TotalResponseTime, server))
		}
	}
	fmt.Println(weights)

	isEqual := allEqual(weights)
	if isEqual {
		numServers := len(servers)
		randServer := rand.IntN(numServers)
		return servers[randServer]
	} else {
		serverToChoose := getMaxWeight(weights)
		return servers[serverToChoose]
	}
}

func calcWeight(alpha float64, totalResponseTime float64, server *server.ServerData) float64 {
	// weight = alpha * (success rate * total response time / server response time)
	// alpha = number of req handled by server / total req received by load balancer
	// alpha = 1 till the total request of the LB becomes 100
	weight := alpha * (server.SuccessRate * totalResponseTime / server.AvgResponseTime)
	// fmt.Printf("The weight of server %d is %.2f\n", server.Id, weight)
	return weight
}

func getMaxWeight(weights []float64) int {
	maxWeight := -99999.0
	serverID := 0

	for id, weight := range weights {
		if weight > maxWeight {
			maxWeight = weight
			serverID = id
		}
	}
	return serverID
}

func allEqual(weight []float64) bool {
	if len(weight) == 0 {
		return true
	}
	first := weight[0]
	// Returns true if no element is NOT equal to the first
	return !slices.ContainsFunc(weight, func(i float64) bool {
		return i != first
	})
}
