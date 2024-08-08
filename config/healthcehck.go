package config

import (
	"fmt"
	"httpFlux/utils"
	"net/http"
)

func CheckHealth(Flux *utils.Flux) bool {
	// Send a health check request to the backend server
	resp, err := http.Get(Flux.URL)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Error in getting response from Flux URL:", err)
		return false
	}
	return true
}
