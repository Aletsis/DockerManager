package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
)

// ListContainers returns list of containers with structured metadata
func (s *Service) ListContainers(ctx context.Context, all bool) ([]ContainerInfo, error) {
	containers, err := s.cli.ContainerList(ctx, container.ListOptions{All: all})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	result := make([]ContainerInfo, 0, len(containers))
	for _, c := range containers {
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		shortID := c.ID
		if len(shortID) > 12 {
			shortID = shortID[:12]
		}

		ports := make([]PortMapping, 0, len(c.Ports))
		for _, p := range c.Ports {
			ports = append(ports, PortMapping{
				IP:          p.IP,
				PrivatePort: p.PrivatePort,
				PublicPort:  p.PublicPort,
				Type:        p.Type,
			})
		}

		composeProject := ""
		composeService := ""
		composeWorkingDir := ""
		composeConfigFile := ""
		if c.Labels != nil {
			composeProject = c.Labels["com.docker.compose.project"]
			composeService = c.Labels["com.docker.compose.service"]
			composeWorkingDir = c.Labels["com.docker.compose.project.working_dir"]
			composeConfigFile = c.Labels["com.docker.compose.project.config_files"]
		}

		networks := make([]ContainerNetworkInfo, 0)
		if c.NetworkSettings != nil && c.NetworkSettings.Networks != nil {
			for netName, netEndpoint := range c.NetworkSettings.Networks {
				if netEndpoint == nil {
					continue
				}
				networks = append(networks, ContainerNetworkInfo{
					NetworkName: netName,
					NetworkID:   netEndpoint.NetworkID,
					IPAddress:   netEndpoint.IPAddress,
					Gateway:     netEndpoint.Gateway,
					MacAddress:  netEndpoint.MacAddress,
					Aliases:     netEndpoint.Aliases,
				})
			}
		}

		result = append(result, ContainerInfo{
			ID:                c.ID,
			ShortID:           shortID,
			Names:             c.Names,
			Name:              name,
			Image:             c.Image,
			ImageID:           c.ImageID,
			Command:           c.Command,
			Created:           c.Created,
			State:             c.State,
			Status:            c.Status,
			Ports:             ports,
			Networks:          networks,
			SizeRw:            c.SizeRw,
			SizeRootFs:        c.SizeRootFs,
			Labels:            c.Labels,
			ComposeProject:    composeProject,
			ComposeService:    composeService,
			ComposeWorkingDir: composeWorkingDir,
			ComposeConfigFile: composeConfigFile,
		})
	}
	return result, nil
}

// StartContainer starts a container by ID
func (s *Service) StartContainer(ctx context.Context, id string) error {
	return s.cli.ContainerStart(ctx, id, container.StartOptions{})
}

// StopContainer stops a container by ID
func (s *Service) StopContainer(ctx context.Context, id string) error {
	timeout := 15
	stopOptions := container.StopOptions{Timeout: &timeout}
	return s.cli.ContainerStop(ctx, id, stopOptions)
}

// RestartContainer restarts a container by ID
func (s *Service) RestartContainer(ctx context.Context, id string) error {
	timeout := 15
	stopOptions := container.StopOptions{Timeout: &timeout}
	return s.cli.ContainerRestart(ctx, id, stopOptions)
}

// PauseContainer pauses a container
func (s *Service) PauseContainer(ctx context.Context, id string) error {
	return s.cli.ContainerPause(ctx, id)
}

// UnpauseContainer unpauses a container
func (s *Service) UnpauseContainer(ctx context.Context, id string) error {
	return s.cli.ContainerUnpause(ctx, id)
}

// RemoveContainer removes a container
func (s *Service) RemoveContainer(ctx context.Context, id string, force bool) error {
	return s.cli.ContainerRemove(ctx, id, container.RemoveOptions{
		Force:         force,
		RemoveVolumes: true,
	})
}

// GetContainerLogs returns recent logs formatted as plain text
func (s *Service) GetContainerLogs(ctx context.Context, id string, tail int) (string, error) {
	tailStr := strconv.Itoa(tail)
	if tail <= 0 {
		tailStr = "200"
	}

	reader, err := s.cli.ContainerLogs(ctx, id, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tailStr,
		Timestamps: false,
	})
	if err != nil {
		return "", fmt.Errorf("failed to fetch container logs: %w", err)
	}
	defer reader.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, reader)
	if err != nil {
		// If raw reader (TTY container enabled), fallback to direct copy
		raw, readErr := io.ReadAll(reader)
		if readErr == nil && len(raw) > 0 {
			return string(raw), nil
		}
		return stdoutBuf.String() + stderrBuf.String(), nil
	}

	var combined strings.Builder
	if stdoutBuf.Len() > 0 {
		combined.WriteString(stdoutBuf.String())
	}
	if stderrBuf.Len() > 0 {
		combined.WriteString(stderrBuf.String())
	}

	return combined.String(), nil
}

// GetContainerStats returns single snapshot metrics for CPU and Memory
func (s *Service) GetContainerStats(ctx context.Context, id string) (*ContainerStats, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	statsResponse, err := s.cli.ContainerStats(ctxTimeout, id, false)
	if err != nil {
		return nil, fmt.Errorf("failed to query stats: %w", err)
	}
	defer statsResponse.Body.Close()

	var stats container.StatsResponse
	decoder := json.NewDecoder(statsResponse.Body)
	if err := decoder.Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode stats: %w", err)
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

	// Networks
	var rx, tx uint64
	for _, net := range stats.Networks {
		rx += net.RxBytes
		tx += net.TxBytes
	}

	// Block I/O
	var readBytes, writeBytes uint64
	for _, bio := range stats.BlkioStats.IoServiceBytesRecursive {
		op := strings.ToLower(bio.Op)
		if op == "read" {
			readBytes += bio.Value
		} else if op == "write" {
			writeBytes += bio.Value
		}
	}

	return &ContainerStats{
		ID:               stats.ID,
		Name:             stats.Name,
		CPUPercentage:    cpuPercent,
		MemoryUsage:      memUsage,
		MemoryLimit:      memLimit,
		MemoryPercentage: memPercent,
		NetworkRx:        rx,
		NetworkTx:        tx,
		BlockRead:        readBytes,
		BlockWrite:       writeBytes,
		PIDs:             stats.PidsStats.Current,
	}, nil
}

// CreateContainer creates and optionally starts a new container based on CreateContainerRequest
func (s *Service) CreateContainer(ctx context.Context, req CreateContainerRequest) (*CreateContainerResult, error) {
	req.Image = strings.TrimSpace(req.Image)
	if req.Image == "" {
		return nil, fmt.Errorf("el nombre de la imagen es requerido")
	}
	req.Name = strings.TrimSpace(req.Name)

	// Verify if image exists locally; if not, pull it automatically
	_, _, err := s.cli.ImageInspectWithRaw(ctx, req.Image)
	if err != nil {
		if pullErr := s.PullImage(ctx, req.Image, nil); pullErr != nil {
			return nil, fmt.Errorf("la imagen %s no está disponible y falló la descarga: %w", req.Image, pullErr)
		}
	}

	// Parse ports if specified
	var exposedPorts nat.PortSet
	var portBindings nat.PortMap
	if len(req.Ports) > 0 {
		var validPorts []string
		for _, p := range req.Ports {
			p = strings.TrimSpace(p)
			if p != "" {
				validPorts = append(validPorts, p)
			}
		}
		if len(validPorts) > 0 {
			var parseErr error
			exposedPorts, portBindings, parseErr = nat.ParsePortSpecs(validPorts)
			if parseErr != nil {
				return nil, fmt.Errorf("error al interpretar especificación de puertos: %w", parseErr)
			}
		}
	}

	// Filter and clean volume binds
	var cleanVolumes []string
	for _, v := range req.Volumes {
		v = strings.TrimSpace(v)
		if v != "" {
			cleanVolumes = append(cleanVolumes, v)
		}
	}

	// Filter and clean environment variables
	var cleanEnv []string
	for _, e := range req.Env {
		e = strings.TrimSpace(e)
		if e != "" {
			cleanEnv = append(cleanEnv, e)
		}
	}

	containerConfig := &container.Config{
		Image:        req.Image,
		Env:          cleanEnv,
		ExposedPorts: exposedPorts,
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		Binds:        cleanVolumes,
	}

	if req.RestartPolicy != "" {
		hostConfig.RestartPolicy = container.RestartPolicy{
			Name: container.RestartPolicyMode(req.RestartPolicy),
		}
	}

	createResp, err := s.cli.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		req.Name,
	)
	if err != nil {
		return nil, fmt.Errorf("error al crear el contenedor: %w", err)
	}

	if req.AutoStart {
		if startErr := s.cli.ContainerStart(ctx, createResp.ID, container.StartOptions{}); startErr != nil {
			return &CreateContainerResult{
				ID:       createResp.ID,
				Warnings: createResp.Warnings,
			}, fmt.Errorf("contenedor creado con ID %s pero falló al iniciar: %w", createResp.ID, startErr)
		}
	}

	return &CreateContainerResult{
		ID:       createResp.ID,
		Warnings: createResp.Warnings,
	}, nil
}
