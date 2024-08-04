package scheduler

import "httpFlux/utils"


func RoundRobbin(index int, backends []*utils.Flux) int {
	return (index + 1) % len(backends)

}
