package loadbalancer

import (
	"fmt"
	"loadBalancer/internal/algorithms"
	"loadBalancer/internal/server"
	"net/http"
	"sync"
	"time"
)

type Data struct {
	Servers    []*server.ServerData
	Algorithm  string
	TotalReqs  int
	FailedReqs int
	StartTime  time.Time
	mu         sync.Mutex
}

var httpClient = &http.Client{}
var balancerData = &Data{}

func handler(w http.ResponseWriter, req *http.Request) {
	balancerData.mu.Lock()
	chosenServer := algorithms.LeastConnections(balancerData.Servers, &balancerData.mu)
	fmt.Println("Request handled by:", chosenServer.Address)
	chosenServer.ActiveConn++
	balancerData.mu.Unlock()

	url := fmt.Sprintf("%s%s", chosenServer.Address, req.URL)
	resp, err := httpClient.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		resp.Body.Close()
		balancerData.mu.Lock()
		chosenServer.ActiveConn--
		balancerData.mu.Unlock()
	}()
	// fmt.Println(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Server did not process", resp.StatusCode)
	}
}

func StartServer(data *Data) error {
	httpClient.Transport = &http.Transport{
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: 1000,
		MaxConnsPerHost:     0,
	}
	balancerData = data
	http.HandleFunc("/getImage", handler)
	http.HandleFunc("/sendVideo", handler)
	fmt.Println("Load balancing Server started at port 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
