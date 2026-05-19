package algorithms

import (
	"fmt"
	"loadBalancer/internal/server"
)

type LeastConnections struct{}

func (lc LeastConnections) NextServer(serversData []*server.ServerData) *server.ServerData {
	var serverToChoose *server.ServerData
	leastConn := 999999
	for _, server := range serversData {
		activeConn := server.ActiveConn
		fmt.Printf("ID:%d, Active conn:%d\n", server.Id, activeConn)
		if activeConn < leastConn {
			leastConn = activeConn
			serverToChoose = server
		}
	}
	return serverToChoose
}
