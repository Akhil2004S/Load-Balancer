package algorithms

import (
	"hash/fnv"
	"loadBalancer/internal/server"
)

func IPHashing(servers []*server.ServerData, clientIP string) *server.ServerData {
	hasher := fnv.New64()
	hasher.Write([]byte(clientIP))
	hashValue := hasher.Sum64()
	serverToChoose := hashValue % uint64(len(servers))
	return servers[serverToChoose]
}
