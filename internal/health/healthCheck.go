package health

import (
	"fmt"
	serverData "loadBalancer/internal/server"
	"net/http"
	"sync"
	"time"
)

func checkHealth(server *serverData.ServerData, activeServers *[]*serverData.ServerData, mu *sync.Mutex) {
	defer func() {
		if r := recover(); r != nil {
			if server.HealthState == serverData.Healthy || server.HealthState == serverData.Evaluating {
				mu.Lock()
				server.HealthState = serverData.Evaluating
				removeServer(server, activeServers)
				server.HealthyCount = 0
				server.UnhealthyCount++
				mu.Unlock()
			}
			if server.HealthState == serverData.Evaluating && (server.UnhealthyCount >= server.UnhealthyThreshold) {
				mu.Lock()
				server.HealthState = serverData.Unhealthy
				server.UnhealthyCount = 0
				removeServer(server, activeServers)
				mu.Unlock()
				fmt.Println("Server marked as unhealthy and removed from list. ID:", server.Id)
			}
		}
	}()
	url := fmt.Sprintf("%s/health", server.Address)
	resp, _ := http.Get(url)
	fmt.Println(activeServers)

	// FIX THIS
	if resp.StatusCode == http.StatusOK {
		mu.Lock()
		if server.HealthState == serverData.Unhealthy || server.HealthState == serverData.Evaluating {
			switch server.HealthState {
			case serverData.Evaluating:
				server.UnhealthyCount = 0
				server.HealthyCount++
			case serverData.Unhealthy:
				server.HealthState = serverData.Evaluating
				removeServer(server, activeServers)
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

	// FIX THIS BLOCK
	if server.HealthState == serverData.Evaluating && (server.HealthyCount >= server.HealthyThreshold) {
		mu.Lock()
		server.HealthState = serverData.Healthy
		server.HealthyCount = 0
		fmt.Printf("Server with ID:%d has passed the healhy threshold", server.Id)
		*activeServers = append(*activeServers, server)
		fmt.Println(activeServers)
		mu.Unlock()
		fmt.Println("Server marked as healthy and added to server list. ID:", server.Id)
	} else if server.HealthState == serverData.Evaluating && (server.UnhealthyCount >= server.UnhealthyThreshold) {
		mu.Lock()
		server.HealthState = serverData.Unhealthy
		server.UnhealthyCount = 0
		removeServer(server, activeServers)
		mu.Unlock()
		fmt.Println("Server marked as unhealthy. ID:", server.Id)
	}
}

func removeServer(serverToRemove *serverData.ServerData, activeServers *[]*serverData.ServerData) {
	var indexToRemove int
	if len(*activeServers) == 0 {
		return
	}
	for id, server := range *activeServers {
		if serverToRemove.Id == server.Id {
			indexToRemove = id
		} else {
			fmt.Println("No server to remove")
		}
	}
	serversList := *activeServers
	fmt.Printf("The index to remove is %d and the active server list is %v\n", indexToRemove, *activeServers)
	serverList := append(serversList[:indexToRemove], serversList[indexToRemove+1:]...)
	*activeServers = serverList
	fmt.Println("Server removed from active server list.", activeServers)
}

func StartHealthChecker(done chan bool, servers []*serverData.ServerData, activeServers *[]*serverData.ServerData, mu *sync.Mutex) {
	fmt.Println("Health Checker started")
	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				for _, server := range servers {
					checkHealth(server, activeServers, mu)
				}
			case <-done:
				fmt.Println("Health checker exited")
				ticker.Stop()
				return
			}
		}
	}()
}
