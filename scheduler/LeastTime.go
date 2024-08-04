package scheduler

import "httpFlux/utils"



func LeastTime(backends []*utils.Flux) int {
	var min int
	min = int(backends[0].ResponseTime)
	index := 0
	for i := 0; i < len(backends); i++ {
		if int(backends[i].ResponseTime) < min {
			min = int(backends[i].ResponseTime)
			index = i
		}
	}
	return index

}
