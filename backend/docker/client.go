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
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// Service encapsulates Docker client operations
type Service struct {
	cli             *client.Client
	terminalManager *TerminalManager
}

// NewService instantiates and connects a Docker client
func NewService() (*Service, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}
	return &Service{
		cli:             cli,
		terminalManager: NewTerminalManager(),
	}, nil
}

// Close closes the Docker client connection and cleans up sessions
func (s *Service) Close() error {
	if s.terminalManager != nil {
		s.terminalManager.CloseAll()
	}
	if s.cli != nil {
		return s.cli.Close()
	}
	return nil
}

// GetOverview returns summary metrics of the Docker host
func (s *Service) GetOverview(ctx context.Context) (*SystemOverview, error) {
	info, err := s.cli.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve docker info: %w", err)
	}

	version, err := s.cli.ServerVersion(ctx)
	versionStr := ""
	if err == nil {
		versionStr = version.Version
	}

	images, _ := s.cli.ImageList(ctx, image.ListOptions{})

	return &SystemOverview{
		Containers:        info.Containers,
		ContainersRunning: info.ContainersRunning,
		ContainersPaused:  info.ContainersPaused,
		ContainersStopped: info.ContainersStopped,
		Images:            len(images),
		ServerVersion:     versionStr,
		OperatingSystem:   info.OperatingSystem,
		NCPU:              info.NCPU,
		MemTotal:          info.MemTotal,
	}, nil
}

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

		result = append(result, ContainerInfo{
			ID:         c.ID,
			ShortID:    shortID,
			Names:      c.Names,
			Name:       name,
			Image:      c.Image,
			ImageID:    c.ImageID,
			Command:    c.Command,
			Created:    c.Created,
			State:      c.State,
			Status:     c.Status,
			Ports:      ports,
			SizeRw:     c.SizeRw,
			SizeRootFs: c.SizeRootFs,
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

// StartTerminal starts an interactive PTY session for a container
func (s *Service) StartTerminal(
	ctx context.Context,
	containerID string,
	shell string,
	rows uint,
	cols uint,
	onData func(sessionID string, chunkBase64 string),
	onExit func(sessionID string),
) (*TerminalStartResult, error) {
	if s.terminalManager == nil {
		return nil, fmt.Errorf("terminal manager not initialized")
	}
	return s.terminalManager.StartSession(ctx, s.cli, containerID, shell, rows, cols, onData, onExit)
}

// WriteTerminal sends input data to a terminal session
func (s *Service) WriteTerminal(sessionID string, data string) error {
	if s.terminalManager == nil {
		return fmt.Errorf("terminal manager not initialized")
	}
	return s.terminalManager.Write(sessionID, data)
}

// ResizeTerminal updates the terminal window dimensions
func (s *Service) ResizeTerminal(ctx context.Context, sessionID string, rows uint, cols uint) error {
	if s.terminalManager == nil {
		return fmt.Errorf("terminal manager not initialized")
	}
	return s.terminalManager.Resize(ctx, s.cli, sessionID, rows, cols)
}

// CloseTerminal terminates an active terminal session
func (s *Service) CloseTerminal(sessionID string) error {
	if s.terminalManager == nil {
		return nil
	}
	s.terminalManager.CloseSession(sessionID)
	return nil
}
