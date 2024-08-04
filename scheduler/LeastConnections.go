package scheduler

import "httpFlux/utils"


func LeastConnections(backends []*utils.Flux) int {
	var min int
	min = int(backends[0].Connections)
	index := 0
	for i := 0; i < len(backends); i++ {
		if int(backends[i].Connections) < min {
			min = int(backends[i].Connections)
		}
	}
	return index

}
