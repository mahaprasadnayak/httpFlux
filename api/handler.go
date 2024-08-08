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

var (
	algos           = utils.WeightedRoundRobbin
	nextServerIndex = 0
	CurrentWeight   = 0
	backendServers  []*utils.Flux
	mutex           sync.Mutex
)

// CustomTransport wraps around the standard Transport to capture errors
type CustomTransport struct {
	Transport http.RoundTripper
}

func (c *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := c.Transport.RoundTrip(req)
	if err != nil {
		// Log the error and return it
		fmt.Printf("Proxy request error: %v\n", err)
	}
	return resp, err
}

func Requesthandler(w http.ResponseWriter, r *http.Request) {
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
	startTime := time.Now()
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.Transport = &CustomTransport{
		Transport: &http.Transport{
			IdleConnTimeout:     90 * time.Second,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
	proxy.ServeHTTP(w, r)
	backendFlux.ResponseTime = time.Since(startTime)
	fmt.Println("Response time:", backendFlux.ResponseTime)

	// Decrement the number of connections for the selected backend server
	defer backendFlux.DecrementConnections()
}

func selectHealthyBackend() *utils.Flux {
	// Filter out unhealthy backends
	healthyBackends := make([]*utils.Flux, 0)
	for _, backend := range backendServers {
		if backend.Healthy {
			healthyBackends = append(healthyBackends, backend)
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
		for _, backend := range backendServers {
			if !config.CheckHealth(backend) {
				// Mark backend as unhealthy
				backend.Healthy = false
				fmt.Printf("Backend %s is unhealthy %d\n", backend.URL, backend.Connections)
			} else {
				// Mark backend as healthy
				backend.Healthy = true
				fmt.Printf("Backend %s is healthy %d\n", backend.URL, backend.Connections)
			}
		}
	}
}

func InitBackendServers(servers []*utils.Flux) {
	backendServers = servers
	go healthCheck()
}
