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

// ComposeDeployRequest holds parameters for deploying a compose stack
type ComposeDeployRequest struct {
	ProjectName   string `json:"projectName"`
	WorkingDir    string `json:"workingDir"`
	ConfigFile    string `json:"configFile"`
	Content       string `json:"content,omitempty"`
	RemoveOrphans bool   `json:"removeOrphans"`
}

// ComposeDownRequest holds parameters for tearing down a compose stack
type ComposeDownRequest struct {
	ProjectName   string `json:"projectName"`
	WorkingDir    string `json:"workingDir"`
	ConfigFile    string `json:"configFile"`
	RemoveVolumes bool   `json:"removeVolumes"`
}

// DockerfileInfo contains information about a discovered Dockerfile
type DockerfileInfo struct {
	Name         string `json:"name"`                   // e.g. "Dockerfile", "backend/Dockerfile"
	Path         string `json:"path"`                   // absolute path on host
	ServiceName  string `json:"serviceName,omitempty"`  // associated compose service (if detected)
	Content      string `json:"content"`                // text content of the Dockerfile
	ExistsOnDisk bool   `json:"existsOnDisk"`
}

// ComposeFileInfo contains the resolved file path and content of a compose file and associated Dockerfiles
type ComposeFileInfo struct {
	ProjectName  string           `json:"projectName"`
	WorkingDir   string           `json:"workingDir"`
	ConfigFile   string           `json:"configFile"`
	Content      string           `json:"content"`
	ExistsOnDisk bool             `json:"existsOnDisk"`
	Dockerfiles  []DockerfileInfo `json:"dockerfiles"`
}


