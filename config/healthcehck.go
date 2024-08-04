package config

import (
	"fmt"
	"httpFlux/utils"
	"net/http"
	"time"
)

func init() {
	// Start health check in the background
	go InithealthCheck()
}


func InithealthCheck() {
	var FluxServer        []*utils.Flux
	for {
		time.Sleep(10 * time.Second)
		// Perform health check for each backend server
		for i := range FluxServer {
			if !CheckHealth(FluxServer[i]) {
				// Mark backend as unhealthy
				FluxServer[i].Healthy = false
				fmt.Printf("Backend %s is unhealthy %d \n", FluxServer[i].URL, FluxServer[i].Connections)
			} else {
				// Mark backend as healthy
				FluxServer[i].Healthy = true
				fmt.Printf("Backend %s is healthy %d \n", FluxServer[i].URL, FluxServer[i].Connections)
			}
		}
	}
}



func CheckHealth(Flux *utils.Flux) bool {
	// Send a health check request to the backend server
	resp, err := http.Get(Flux.URL)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Error in Getting Response from Flux Url")
		return false
	}
	return true
}
