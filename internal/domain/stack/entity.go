package stack

import "dockermanager/internal/domain/container"

// ComposeStack is a domain aggregate grouping containers by compose project
type ComposeStack struct {
	Name         string                `json:"name"`
	WorkingDir   string                `json:"workingDir,omitempty"`
	ConfigFile   string                `json:"configFile,omitempty"`
	Containers   []container.Container `json:"containers"`
	RunningCount int                   `json:"runningCount"`
	TotalCount   int                   `json:"totalCount"`
}

// RecalculateCounts recalculates total and running containers
func (s *ComposeStack) RecalculateCounts() {
	s.TotalCount = len(s.Containers)
	s.RunningCount = 0
	for _, c := range s.Containers {
		if c.IsRunning() {
			s.RunningCount++
		}
	}
}
