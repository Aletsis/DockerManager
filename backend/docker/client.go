package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
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

// StartStack starts all non-running containers belonging to the specified compose project
func (s *Service) StartStack(ctx context.Context, projectName string) error {
	projectName = strings.TrimSpace(projectName)
	if projectName == "" {
		return fmt.Errorf("el nombre del proyecto compose no puede estar vacío")
	}

	filterArgs := filters.NewArgs()
	filterArgs.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))

	containers, err := s.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filterArgs})
	if err != nil {
		return fmt.Errorf("error al listar contenedores del stack %s: %w", projectName, err)
	}
	if len(containers) == 0 {
		return fmt.Errorf("no se encontraron contenedores para el stack %q", projectName)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []string

	for _, c := range containers {
		if c.State == "running" {
			continue
		}
		wg.Add(1)
		go func(cid string) {
			defer wg.Done()
			if err := s.StartContainer(ctx, cid); err != nil {
				mu.Lock()
				short := cid
				if len(short) > 12 {
					short = short[:12]
				}
				errs = append(errs, fmt.Sprintf("%s: %v", short, err))
				mu.Unlock()
			}
		}(c.ID)
	}
	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("errores al iniciar el stack: %s", strings.Join(errs, "; "))
	}
	return nil
}

// StopStack stops all running containers belonging to the specified compose project
func (s *Service) StopStack(ctx context.Context, projectName string) error {
	projectName = strings.TrimSpace(projectName)
	if projectName == "" {
		return fmt.Errorf("el nombre del proyecto compose no puede estar vacío")
	}

	filterArgs := filters.NewArgs()
	filterArgs.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))

	containers, err := s.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filterArgs})
	if err != nil {
		return fmt.Errorf("error al listar contenedores del stack %s: %w", projectName, err)
	}
	if len(containers) == 0 {
		return fmt.Errorf("no se encontraron contenedores para el stack %q", projectName)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []string

	for _, c := range containers {
		if c.State != "running" {
			continue
		}
		wg.Add(1)
		go func(cid string) {
			defer wg.Done()
			if err := s.StopContainer(ctx, cid); err != nil {
				mu.Lock()
				short := cid
				if len(short) > 12 {
					short = short[:12]
				}
				errs = append(errs, fmt.Sprintf("%s: %v", short, err))
				mu.Unlock()
			}
		}(c.ID)
	}
	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("errores al detener el stack: %s", strings.Join(errs, "; "))
	}
	return nil
}

// RestartStack restarts all containers belonging to the specified compose project
func (s *Service) RestartStack(ctx context.Context, projectName string) error {
	projectName = strings.TrimSpace(projectName)
	if projectName == "" {
		return fmt.Errorf("el nombre del proyecto compose no puede estar vacío")
	}

	filterArgs := filters.NewArgs()
	filterArgs.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))

	containers, err := s.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filterArgs})
	if err != nil {
		return fmt.Errorf("error al listar contenedores del stack %s: %w", projectName, err)
	}
	if len(containers) == 0 {
		return fmt.Errorf("no se encontraron contenedores para el stack %q", projectName)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []string

	for _, c := range containers {
		wg.Add(1)
		go func(cid string) {
			defer wg.Done()
			if err := s.RestartContainer(ctx, cid); err != nil {
				mu.Lock()
				short := cid
				if len(short) > 12 {
					short = short[:12]
				}
				errs = append(errs, fmt.Sprintf("%s: %v", short, err))
				mu.Unlock()
			}
		}(c.ID)
	}
	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("errores al reiniciar el stack: %s", strings.Join(errs, "; "))
	}
	return nil
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

// ListImages returns list of all local images with structured metadata and usage status
func (s *Service) ListImages(ctx context.Context) ([]ImageInfo, error) {
	images, err := s.cli.ImageList(ctx, image.ListOptions{All: false})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	// Cross reference running/stopped containers to find images in use
	containers, _ := s.cli.ContainerList(ctx, container.ListOptions{All: true})
	usedImages := make(map[string]bool)
	for _, c := range containers {
		usedImages[c.ImageID] = true
		usedImages[c.Image] = true
	}

	result := make([]ImageInfo, 0, len(images))
	for _, img := range images {
		shortID := img.ID
		if strings.HasPrefix(shortID, "sha256:") {
			shortID = strings.TrimPrefix(shortID, "sha256:")
		}
		if len(shortID) > 12 {
			shortID = shortID[:12]
		}

		isDangling := false
		repo := "<none>"
		tag := "<none>"

		if len(img.RepoTags) > 0 && img.RepoTags[0] != "<none>:<none>" {
			parts := strings.Split(img.RepoTags[0], ":")
			if len(parts) >= 2 {
				repo = strings.Join(parts[:len(parts)-1], ":")
				tag = parts[len(parts)-1]
			} else {
				repo = img.RepoTags[0]
				tag = "latest"
			}
		} else {
			isDangling = true
		}

		inUse := img.Containers > 0 || usedImages[img.ID]
		if !inUse {
			for _, rt := range img.RepoTags {
				if usedImages[rt] {
					inUse = true
					break
				}
			}
		}

		result = append(result, ImageInfo{
			ID:         img.ID,
			ShortID:    shortID,
			Repository: repo,
			Tag:        tag,
			RepoTags:   img.RepoTags,
			Created:    img.Created,
			Size:       img.Size,
			SharedSize: img.SharedSize,
			Containers: img.Containers,
			InUse:      inUse,
			IsDangling: isDangling,
		})
	}
	return result, nil
}

// GetDiskUsage returns aggregated metrics of image disk storage and recoverable space
func (s *Service) GetDiskUsage(ctx context.Context) (*DiskUsageSummary, error) {
	du, err := s.cli.DiskUsage(ctx, types.DiskUsageOptions{})
	if err != nil {
		// Fallback: calculate directly from ListImages if DiskUsage endpoint is restricted
		images, listErr := s.ListImages(ctx)
		if listErr != nil {
			return nil, fmt.Errorf("failed to retrieve disk usage: %w", err)
		}
		var totalSize, danglingSize, reclaimableSize int64
		danglingCount := 0
		for _, img := range images {
			totalSize += img.Size
			if img.IsDangling {
				danglingCount++
				danglingSize += img.Size
			}
			if !img.InUse {
				reclaimableSize += img.Size
			}
		}
		return &DiskUsageSummary{
			TotalImages:     len(images),
			TotalSize:       totalSize,
			DanglingCount:   danglingCount,
			DanglingSize:    danglingSize,
			ReclaimableSize: reclaimableSize,
		}, nil
	}

	var totalSize, danglingSize, reclaimableSize int64
	danglingCount := 0

	containers, _ := s.cli.ContainerList(ctx, container.ListOptions{All: true})
	usedImages := make(map[string]bool)
	for _, c := range containers {
		usedImages[c.ImageID] = true
		usedImages[c.Image] = true
	}

	for _, img := range du.Images {
		totalSize += img.Size
		isDangling := len(img.RepoTags) == 0 || (len(img.RepoTags) == 1 && img.RepoTags[0] == "<none>:<none>")
		if isDangling {
			danglingCount++
			danglingSize += img.Size
		}
		inUse := img.Containers > 0 || usedImages[img.ID]
		if !inUse {
			for _, rt := range img.RepoTags {
				if usedImages[rt] {
					inUse = true
					break
				}
			}
		}
		if !inUse {
			reclaimableSize += img.Size
		}
	}

	return &DiskUsageSummary{
		TotalImages:     len(du.Images),
		TotalSize:       totalSize,
		DanglingCount:   danglingCount,
		DanglingSize:    danglingSize,
		ReclaimableSize: reclaimableSize,
	}, nil
}

// PullImage downloads an image while streaming progress events
func (s *Service) PullImage(ctx context.Context, ref string, onProgress func(event PullProgressEvent)) error {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return fmt.Errorf("el nombre de la imagen no puede estar vacío")
	}

	reader, err := s.cli.ImagePull(ctx, ref, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("error al iniciar descarga de imagen %s: %w", ref, err)
	}
	defer reader.Close()

	decoder := json.NewDecoder(reader)
	for {
		var msg jsonmessage.JSONMessage
		if err := decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error al leer respuesta de descarga: %w", err)
		}

		if msg.Error != nil {
			if onProgress != nil {
				onProgress(PullProgressEvent{
					ID:     msg.ID,
					Status: "Error",
					Error:  msg.Error.Message,
				})
			}
			return fmt.Errorf("error descargando imagen: %s", msg.Error.Message)
		}

		if onProgress != nil {
			var cur, tot int64
			if msg.Progress != nil {
				cur = msg.Progress.Current
				tot = msg.Progress.Total
			}
			onProgress(PullProgressEvent{
				ID:       msg.ID,
				Status:   msg.Status,
				Progress: msg.ProgressMessage,
				Current:  cur,
				Total:    tot,
			})
		}
	}

	return nil
}

// RemoveImage deletes an image by ID or name
func (s *Service) RemoveImage(ctx context.Context, id string, force bool) error {
	_, err := s.cli.ImageRemove(ctx, id, image.RemoveOptions{
		Force:         force,
		PruneChildren: true,
	})
	if err != nil {
		return fmt.Errorf("error al eliminar imagen %s: %w", id, err)
	}
	return nil
}

// PruneImages deletes unused or dangling images and returns reclaimed space
func (s *Service) PruneImages(ctx context.Context, danglingOnly bool) (*PruneResult, error) {
	pruneFilters := filters.NewArgs()
	if danglingOnly {
		pruneFilters.Add("dangling", "true")
	} else {
		pruneFilters.Add("dangling", "false")
	}

	report, err := s.cli.ImagesPrune(ctx, pruneFilters)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar limpieza de imágenes: %w", err)
	}

	deletedIDs := make([]string, 0, len(report.ImagesDeleted))
	for _, item := range report.ImagesDeleted {
		if item.Deleted != "" {
			deletedIDs = append(deletedIDs, item.Deleted)
		} else if item.Untagged != "" {
			deletedIDs = append(deletedIDs, item.Untagged)
		}
	}

	return &PruneResult{
		ImagesDeleted:  deletedIDs,
		SpaceReclaimed: report.SpaceReclaimed,
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

