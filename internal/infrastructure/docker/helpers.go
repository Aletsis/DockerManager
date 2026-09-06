package docker

import (
	"fmt"
	"strings"

	containerdomain "dockermanager/internal/domain/container"

	dockertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

// ParseRepoTag parses a raw image tag string into repository, tag, and dangling indicator
func ParseRepoTag(rawTag string) (repo string, tag string, isDangling bool) {
	rawTag = strings.TrimSpace(rawTag)
	if rawTag == "" || rawTag == "<none>:<none>" || rawTag == "<none>" {
		return "<none>", "<none>", true
	}

	parts := strings.Split(rawTag, ":")
	if len(parts) >= 2 {
		repo = strings.Join(parts[:len(parts)-1], ":")
		tag = parts[len(parts)-1]
		return repo, tag, false
	}

	return rawTag, "latest", false
}

// CalculateContainerStats extracts and calculates normalized resource utilization from Docker stats
func CalculateContainerStats(id string, name string, stats *dockertypes.StatsResponse) *containerdomain.Stats {
	if stats == nil {
		return &containerdomain.Stats{ID: id, Name: name}
	}

	// Calculate CPU Percentage
	cpuPercent := 0.0
	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)
	onlineCPUs := float64(stats.CPUStats.OnlineCPUs)
	if onlineCPUs == 0 {
		onlineCPUs = float64(len(stats.CPUStats.CPUUsage.PercpuUsage))
	}
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}

	// Calculate Memory Usage
	memUsage := stats.MemoryStats.Usage
	if cache, ok := stats.MemoryStats.Stats["inactive_file"]; ok {
		if memUsage > cache {
			memUsage -= cache
		}
	} else if cache, ok := stats.MemoryStats.Stats["cache"]; ok {
		if memUsage > cache {
			memUsage -= cache
		}
	}

	memLimit := stats.MemoryStats.Limit
	memPercent := 0.0
	if memLimit > 0 {
		memPercent = (float64(memUsage) / float64(memLimit)) * 100.0
	}

	// Calculate Network I/O
	var rx, tx uint64
	for _, net := range stats.Networks {
		rx += net.RxBytes
		tx += net.TxBytes
	}

	// Calculate Block I/O
	var readBytes, writeBytes uint64
	for _, bio := range stats.BlkioStats.IoServiceBytesRecursive {
		op := strings.ToLower(bio.Op)
		if op == "read" {
			readBytes += bio.Value
		} else if op == "write" {
			writeBytes += bio.Value
		}
	}

	return &containerdomain.Stats{
		ID:               id,
		Name:             name,
		CPUPercentage:    cpuPercent,
		MemoryUsage:      memUsage,
		MemoryLimit:      memLimit,
		MemoryPercentage: memPercent,
		NetworkRx:        rx,
		NetworkTx:        tx,
		BlockRead:        readBytes,
		BlockWrite:       writeBytes,
		Pids:             stats.PidsStats.Current,
	}
}

// BuildContainerConfig validates and constructs container Config and HostConfig
func BuildContainerConfig(spec containerdomain.CreateSpec) (*dockertypes.Config, *dockertypes.HostConfig, error) {
	spec.Image = strings.TrimSpace(spec.Image)
	if spec.Image == "" {
		return nil, nil, fmt.Errorf("el nombre de la imagen es requerido")
	}

	// Parse ports if specified
	var exposedPorts nat.PortSet
	var portBindings nat.PortMap
	if len(spec.Ports) > 0 {
		var validPorts []string
		for _, p := range spec.Ports {
			p = strings.TrimSpace(p)
			if p != "" {
				validPorts = append(validPorts, p)
			}
		}
		if len(validPorts) > 0 {
			var parseErr error
			exposedPorts, portBindings, parseErr = nat.ParsePortSpecs(validPorts)
			if parseErr != nil {
				return nil, nil, fmt.Errorf("error al interpretar especificación de puertos: %w", parseErr)
			}
		}
	}

	// Filter and clean volume binds
	var cleanVolumes []string
	for _, v := range spec.Volumes {
		v = strings.TrimSpace(v)
		if v != "" {
			cleanVolumes = append(cleanVolumes, v)
		}
	}

	// Filter and clean environment variables
	var cleanEnv []string
	for _, e := range spec.Env {
		e = strings.TrimSpace(e)
		if e != "" {
			cleanEnv = append(cleanEnv, e)
		}
	}

	containerConfig := &dockertypes.Config{
		Image:        spec.Image,
		Cmd:          spec.Cmd,
		Env:          cleanEnv,
		ExposedPorts: exposedPorts,
	}

	hostConfig := &dockertypes.HostConfig{
		PortBindings: portBindings,
		Binds:        cleanVolumes,
	}

	if spec.RestartPolicy != "" {
		hostConfig.RestartPolicy = dockertypes.RestartPolicy{
			Name: dockertypes.RestartPolicyMode(spec.RestartPolicy),
		}
	}

	return containerConfig, hostConfig, nil
}
