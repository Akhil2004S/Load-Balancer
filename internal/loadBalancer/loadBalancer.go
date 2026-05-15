package loadbalancer

import (
	"fmt"
	"loadBalancer/internal/server"
	"net/http"
	"sync"
	"time"
)

type Data struct {
	servers    *[]server.ServerData
	algorithm  string
	TotalReqs  int
	FailedReqs int
	StartTime  time.Time
	mu         sync.Mutex
}

var httpClient = &http.Client{}

func imageHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req)
	resp, err := httpClient.Get("http://127.0.0.1:8080/getImage")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(resp.StatusCode)
}

func videoHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req)
	httpClient.Get("http://127.0.0.1:8080/sendVideo")
}

func StartServer() error {
	http.HandleFunc("/getImage", imageHandler)
	http.HandleFunc("/sendVideo", videoHandler)
	fmt.Println("Load balancing Server started at port 3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
