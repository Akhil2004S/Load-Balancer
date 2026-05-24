package health

import (
	"fmt"
	serverData "loadBalancer/internal/server"
	"net/http"
	"sync"
	"time"
)

func checkHealth(server *serverData.ServerData, activeServers *[]*serverData.ServerData, mu *sync.RWMutex, isStarted chan bool) {
	url := fmt.Sprintf("%s/health", server.Address)
	resp, err := http.Get(url)
	if err != nil {
		mu.Lock()
		if server.HealthState == serverData.Healthy || server.HealthState == serverData.Evaluating {
			switch server.HealthState {
			case serverData.Evaluating:
				server.HealthyCount = 0
				server.UnhealthyCount++
			case serverData.Healthy:
				server.HealthState = serverData.Evaluating
				removeServer(server, activeServers)
				server.UnhealthyCount++
			}
		}
		if server.HealthState == serverData.Evaluating && (server.UnhealthyCount >= server.UnhealthyThreshold) {
			server.HealthState = serverData.Unhealthy
			server.UnhealthyCount = 0
			removeServer(server, activeServers)
			fmt.Println("Server marked as unhealthy. ID:", server.Id)
		}
		mu.Unlock()
		return
	}

	if resp.StatusCode == http.StatusOK {
		mu.Lock()
		if server.HealthState == serverData.Unhealthy || server.HealthState == serverData.Evaluating {
			switch server.HealthState {
			case serverData.Evaluating:
				server.UnhealthyCount = 0
				server.HealthyCount++
			case serverData.Unhealthy:
				server.HealthState = serverData.Evaluating
				server.HealthyCount++
			}
		}
		mu.Unlock()
	} else {
		mu.Lock()
		if server.HealthState == serverData.Healthy || server.HealthState == serverData.Evaluating {
			switch server.HealthState {
			case serverData.Evaluating:
				server.HealthyCount = 0
				server.UnhealthyCount++
			case serverData.Healthy:
				server.HealthState = serverData.Evaluating
				removeServer(server, activeServers)
			}
		}
		mu.Unlock()
	}

	mu.Lock()
	if server.HealthState == serverData.Evaluating && (server.HealthyCount >= server.HealthyThreshold) {
		server.HealthState = serverData.Healthy
		server.HealthyCount = 0
		// fmt.Printf("Server with ID:%d has passed the healhy threshold", server.Id)
		if len(*activeServers) == 0 {
			*activeServers = append(*activeServers, server)
			select {
			case isStarted <- true:
			default:
			}
		} else {
			*activeServers = append(*activeServers, server)
		}
		// fmt.Println(activeServers)
		// fmt.Println("Server marked as healthy and added to server list. ID:", server.Id)
	} else if server.HealthState == serverData.Evaluating && (server.UnhealthyCount >= server.UnhealthyThreshold) {
		server.HealthState = serverData.Unhealthy
		server.UnhealthyCount = 0
		removeServer(server, activeServers)
		fmt.Println("Server marked as unhealthy. ID:", server.Id)
	}
	mu.Unlock()
}

func removeServer(serverToRemove *serverData.ServerData, activeServers *[]*serverData.ServerData) {
	var newServersList []*serverData.ServerData
	if len(*activeServers) == 0 {
		return
	}

	for _, server := range *activeServers {
		if serverToRemove.Id == server.Id {
			continue
		} else {
			newServersList = append(newServersList, server)
		}
	}
	*activeServers = newServersList
	// fmt.Println("Server removed from active server list.", activeServers)
}

func StartHealthChecker(done chan bool, isStarted chan bool, servers []*serverData.ServerData, activeServers *[]*serverData.ServerData, mu *sync.RWMutex) {
	fmt.Println("Health Checker started")
	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				for _, server := range servers {
					checkHealth(server, activeServers, mu, isStarted)
				}
			case <-done:
				fmt.Println("Health checker exited")
				ticker.Stop()
				return
			}
		}
	}()
}
