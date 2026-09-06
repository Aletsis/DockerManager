package network

import "dockermanager/internal/domain/container"

// Group represents a Docker network grouping connected containers
type Group struct {
	Name         string                `json:"name"`
	NetworkID    string                `json:"networkId,omitempty"`
	IsDefault    bool                  `json:"isDefault,omitempty"`
	Containers   []container.Container `json:"containers"`
	RunningCount int                   `json:"runningCount"`
	TotalCount   int                   `json:"totalCount"`
}

// RecalculateCounts updates running and total container tallies
func (g *Group) RecalculateCounts() {
	g.TotalCount = len(g.Containers)
	g.RunningCount = 0
	for _, c := range g.Containers {
		if c.IsRunning() {
			g.RunningCount++
		}
	}
}
