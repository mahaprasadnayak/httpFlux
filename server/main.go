package main

import (
	"fmt"
	"httpFlux/utils"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
)
var (
	mutex           sync.Mutex
	FluxServer        []*utils.Flux
	nextServerIndex = 0
	CurrentWeight   = 0
	algos = utils.WeightedRoundRobbin
) 
func main() {
	cfg, err := utils.FetchFluxConfig("config/config.json")
	if err != nil {
		fmt.Println("Error reading config file", err)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()
	urlString := cfg.ProxyUrl
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}
	_, port, err := net.SplitHostPort(parsedURL.Host)
	strPort := fmt.Sprintf(":%s", port)
	if err != nil {
		fmt.Println("Error parsing port:", err)
		return

	}
	FluxServer = make([]*utils.Flux, len(cfg.Nodes))
	for i, s := range cfg.Nodes {
		FluxServer[i] = &utils.Flux{
			URL:         s.NodeURL,
			Weight:      s.Weight,
			Healthy:     true,
			Connections: 0,
		}
		fmt.Println(" Starting Flux servers", FluxServer[i])
	}
	//http.HandleFunc("/", handler)
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Remote Server status", r.RemoteAddr)
		fmt.Fprintln(w, "HttpFlux is running successfully OK!!!",)
	})
	fmt.Printf("HttpFlux started on %s", port)
	log.Fatal(http.ListenAndServe(strPort, nil))
}
