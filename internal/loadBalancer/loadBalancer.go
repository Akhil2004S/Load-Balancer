package loadbalancer

import (
	"fmt"
	"loadBalancer/internal/algorithms"
	"loadBalancer/internal/health"
	"loadBalancer/internal/server"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Data struct {
	Servers       []*server.ServerData
	ActiveServers []*server.ServerData
	Algorithm     algorithms.Algorithm
	TotalReqs     int
	FailedReqs    int
	Alpha         float64
	StartTime     time.Time
	Done          chan bool
	IsStarted     chan bool
	mu            sync.Mutex
}

var httpClient = &http.Client{}
var balancerData = &Data{}

func handler(w http.ResponseWriter, req *http.Request) {

	balancerData.mu.Lock()

	clientIP, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	var totalResponseTime float64
	for _, server := range balancerData.Servers {
		totalResponseTime += server.AvgResponseTime
	}
	balancerData.TotalReqs++

	switch algo := balancerData.Algorithm.(type) {
	case *algorithms.IPHash:
		algo.ClientIP = clientIP
	case *algorithms.LeastResponseTime:
		algo.TotalRequests = balancerData.TotalReqs
		algo.TotalResponseTime = totalResponseTime
	}

	if len(balancerData.ActiveServers) == 0 {
		balancerData.mu.Unlock()
		http.Error(w, "Servers are not ready yet", http.StatusServiceUnavailable)
		return
	}

	chosenServer := balancerData.Algorithm.NextServer(balancerData.ActiveServers)
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

func StartServer(data *Data, port int) error {
	balancerData = data
	balancerData.StartTime = time.Now()
	balancerData.Done = make(chan bool)
	balancerData.IsStarted = make(chan bool)

	health.StartHealthChecker(balancerData.Done, balancerData.IsStarted, balancerData.Servers, &balancerData.ActiveServers, &balancerData.mu)

	addr := fmt.Sprintf(":%d", port)
	httpServer := &http.Server{
		Addr: addr,
	}

	http.HandleFunc("/getImage", handler)
	http.HandleFunc("/sendVideo", handler)

	<-balancerData.IsStarted
	go func() {
		fmt.Println("Load balancing Server started at port", port)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Listen and serve error:", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	balancerData.Done <- true
	log.Println("Shutting down load balancer")
	fmt.Println("Exiting...")

	return nil
}
