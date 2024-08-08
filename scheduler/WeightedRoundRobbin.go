package scheduler

import (
	"httpFlux/utils"
)

func WeightedRoundRobbin(current int, backends []*utils.Flux) (int, int) {
	totalWeight := 0
	for _, backend := range backends {
		totalWeight += backend.Weight
	}
	current = (current + 1) % totalWeight
	selected := current
	for i, backend := range backends {
		if selected < backend.Weight {
			return i, current
		}
		selected -= backend.Weight
	}
	return 0, 0
}
