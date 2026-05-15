package main

import (
	"fmt"
	loadbalancer "loadBalancer/internal/loadBalancer"
	"loadBalancer/internal/server"
	"log"
)

func main() {
	server := server.CreateServer("127.0.0.1:8080")
	fmt.Println("Server created at port 8080", server)

	if err := loadbalancer.StartServer(); err != nil {
		log.Fatal(err)
	}
	// Should add CLI
}
