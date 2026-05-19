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
	Alpha      float64
	StartTime  time.Time
	mu         sync.Mutex
}

var httpClient = &http.Client{}
var balancerData = &Data{}

func handler(w http.ResponseWriter, req *http.Request) {
	// algorithms.LeastResponseTime(balancerData.Servers)
	balancerData.mu.Lock()
	var totalResponseTime float64
	for _, server := range balancerData.Servers {
		fmt.Printf("Requests handled by server %d is %d\n", server.Id, server.TotalReqs)
		totalResponseTime += server.AvgResponseTime
	}
	balancerData.TotalReqs++
	chosenServer := algorithms.LeastResponseTime(balancerData.Servers, balancerData.TotalReqs, totalResponseTime, &balancerData.mu)
	fmt.Println("Request handled by:", chosenServer.Address)
	balancerData.mu.Unlock()

	url := fmt.Sprintf("%s%s", chosenServer.Address, req.URL)
	balancerData.mu.Lock()
	chosenServer.ActiveConn++
	chosenServer.TotalReqs++
	balancerData.mu.Unlock()
	startTime := time.Now()
	resp, err := httpClient.Get(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		balancerData.mu.Lock()
		chosenServer.SuccessReqs++
		if chosenServer.AvgResponseTime == 0 {
			chosenServer.AvgResponseTime = float64(time.Since(startTime).Milliseconds())
		} else {
			// Avg response time using EMA
			// AvgRespTime = alpha * newSample + (1 - alpha) * oldAvg
			chosenServer.AvgResponseTime =
				balancerData.Alpha*float64(time.Since(startTime).Milliseconds()) + (1-balancerData.Alpha)*chosenServer.AvgResponseTime
		}
		chosenServer.SuccessRate = float64(chosenServer.SuccessReqs) / float64(chosenServer.TotalReqs)
		balancerData.mu.Unlock()
	} else {
		http.Error(w, "Server did not process", resp.StatusCode)
	}

	defer func() {
		resp.Body.Close()
		balancerData.mu.Lock()
		chosenServer.ActiveConn--
		balancerData.mu.Unlock()
	}()
}

func StartServer(data *Data) error {
	balancerData = data
	balancerData.StartTime = time.Now()
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
