package main

import (
	"fmt"
	"httpFlux/api"
	"httpFlux/utils"
	"log"
	"net"
	"net/http"
	"net/url"
)

var (
	FluxServer []*utils.Flux
)

func main() {
	cfg, err := utils.FetchFluxConfig("./config/config.json")
	if err != nil {
		fmt.Println("Error reading config file:", err)
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
	if err != nil {
		fmt.Println("Error parsing port:", err)
		return
	}
	strPort := fmt.Sprintf(":%s", port)

	// Initialize FluxServer with the configuration
	FluxServer = make([]*utils.Flux, len(cfg.Nodes))
	for i, s := range cfg.Nodes {
		FluxServer[i] = &utils.Flux{
			URL:         s.NodeURL,
			Weight:      s.Weight,
			Healthy:     true,
			Connections: 0,
		}
		fmt.Println("Starting Flux server", FluxServer[i])
	}

	api.InitBackendServers(FluxServer)

	http.HandleFunc("/", api.Requesthandler)
	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Remote Server status:", r.RemoteAddr)
		fmt.Fprintln(w, "HttpFlux is running successfully OK!!!")
	})
	fmt.Printf("HttpFlux started on %s\n", strPort)
	log.Fatal(http.ListenAndServe(strPort, nil))
}
