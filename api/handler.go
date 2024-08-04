package api

import (
	"fmt"
	"httpFlux/config"
	"httpFlux/scheduler"
	"httpFlux/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)
var(
	algos = utils.WeightedRoundRobbin
	nextServerIndex = 0
	CurrentWeight   = 0
	backends        []*utils.Flux
)
func Requesthandler(w http.ResponseWriter, r *http.Request) {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	backendFlux := selectHealthyBackend()
	if backendFlux == nil {
		fmt.Println("No healthy backend servers available")
		http.Error(w, "No healthy backend servers available", http.StatusServiceUnavailable)
		return
	}
	fmt.Println("Selected backend server:", backendFlux.URL)
	url, _ := url.Parse(backendFlux.URL)

	// Increment the number of connections for the selected backend server

	backendFlux.IncrementConnections()

	// Forward request to selected backend server
	//the below is used to implemnet the keep alive connection upto 100
	startTime := time.Now()
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = &http.Transport{
		IdleConnTimeout:     90 * time.Second,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	proxy.ServeHTTP(w, r)
	backendFlux.ResponseTime = time.Since(startTime)
	fmt.Println("Response time", backendFlux.ResponseTime)

	// Decrement the number of connections for the selected backend server
	defer backendFlux.DecrementConnections()

}


func selectHealthyBackend() *utils.Flux {
	 
	// Filter out unhealthy backends
	healthyBackends := make([]*utils.Flux, 0)
	for i := range backends {
		if backends[i].Healthy {
			healthyBackends = append(healthyBackends, backends[i])
		}
	}

	if len(healthyBackends) == 0 {
		return nil
	}

	switch algos {
	case utils.RoundRobbin:
		nextServerIndex = scheduler.RoundRobbin(nextServerIndex, healthyBackends)
		return healthyBackends[nextServerIndex]
	case utils.LeastConnections:
		nextServerIndex = scheduler.LeastConnections(healthyBackends)
		return healthyBackends[nextServerIndex]
	case utils.WeightedRoundRobbin:
		// fmt.Println("CurrentWeight is the", CurrentWeight)
		nextServerIndex, CurrentWeight = scheduler.WeightedRoundRobbin(CurrentWeight, healthyBackends)
		return healthyBackends[nextServerIndex]
	case utils.LeastTime:
		nextServerIndex = scheduler.LeastTime(healthyBackends)
		return healthyBackends[nextServerIndex]

	default:
		// Use round-robin as the default algorithm
		nextServerIndex = scheduler.RoundRobbin(nextServerIndex, healthyBackends)
		return healthyBackends[nextServerIndex]

	}

}

func healthCheck() {
	for {
		time.Sleep(10 * time.Second) // Check health every 10 seconds
		// Perform health check for each backend server
		for i := range backends {
			if !config.CheckHealth(backends[i]) {
				// Mark backend as unhealthy
				backends[i].Healthy = false
				fmt.Printf("Backend %s is unhealthy %d \n", backends[i].URL, backends[i].Connections)
			} else {
				// Mark backend as healthy
				backends[i].Healthy = true
				fmt.Printf("Backend %s is healthy %d \n", backends[i].URL, backends[i].Connections)
			}
		}
	}
}

func init() {
	// Start health check in the background
	go healthCheck()
}
