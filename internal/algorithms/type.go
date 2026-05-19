package algorithms

import "loadBalancer/internal/server"

type Algorithm interface {
	NextServer([]*server.ServerData) *server.ServerData
}
