package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type FluxConfig struct {
	Nodes       []ServerConfig `json:"nodes"`
	ProxyUrl    string         `json:"proxy"`
}

type ServerConfig struct {
	NodeURL string `json:"node_url"`
	Weight  int    `json:"weight"`
}

type Flux struct {
	URL           string        `json:"url"`
	Healthy       bool
	Connections   int
	ResponseTime  time.Duration
	Weight        int `json:"weight"`
	WeightedScore float64
}

func (b *Flux) IncrementConnections() {
	b.Connections++
	fmt.Println("Connections increment:", b.Connections)
}

func (b *Flux) DecrementConnections() {
	fmt.Println("Connections decrement:", b.Connections)
	b.Connections--
}

var (
	RoundRobbin         = "RoundRobbin"
	LeastConnections    = "LeastConnections"
	LeastTime           = "LeastTime"
	WeightedRoundRobbin = "WeightedRoundRobbin"
)

func FetchFluxConfig(filename string) (*FluxConfig, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config FluxConfig
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
