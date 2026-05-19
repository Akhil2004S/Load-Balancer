package algorithms

import (
	"fmt"
	"hash/fnv"
	"loadBalancer/internal/server"
)

type IPHash struct {
	ClientIP string
}

func (ipHash IPHash) NextServer(servers []*server.ServerData) *server.ServerData {
	fmt.Println(ipHash.ClientIP)
	hasher := fnv.New64()
	hasher.Write([]byte(ipHash.ClientIP))
	hashValue := hasher.Sum64()
	serverToChoose := hashValue % uint64(len(servers))
	return servers[serverToChoose]
}
