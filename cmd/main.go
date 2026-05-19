package main

import (
	"fmt"
	"loadBalancer/internal/algorithms"
	loadbalancer "loadBalancer/internal/loadBalancer"
	"loadBalancer/internal/server"
	"log"
)

func main() {
	server1 := server.CreateServer(1, "http://127.0.0.1:8080")
	fmt.Println("Server created at port 8080")

	server2 := server.CreateServer(2, "http://127.0.0.1:8081")
	fmt.Println("Server created at port 8081")

	data := &loadbalancer.Data{}
	data.Servers = append(data.Servers, server1, server2)
	data.Alpha = 0.1
	data.Algorithm = algorithms.LeastConnections{}
	// data.Algorithm = algorithms.LeastResponseTime{}
	// data.Algorithm = algorithms.IPHash{}
	if err := loadbalancer.StartServer(data); err != nil {
		log.Fatal(err)
	}
	// Should add CLI
}
