package main

import (
	"context"
	"fmt"
	"sync"

	"dockermanager/backend/docker"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct manages application state and exposes Docker services to Wails
type App struct {
	ctx           context.Context
	mu            sync.Mutex
	dockerService *docker.Service
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	svc, err := docker.NewService()
	if err != nil {
		fmt.Println("Warning: Failed to initialize Docker client:", err)
	} else {
		a.dockerService = svc
	}
}

// shutdown is called when the application terminates
func (a *App) shutdown(ctx context.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.dockerService != nil {
		_ = a.dockerService.Close()
		a.dockerService = nil
	}
}

// getService ensures a thread-safe connection to the Docker daemon
func (a *App) getService() (*docker.Service, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.dockerService == nil {
		svc, err := docker.NewService()
		if err != nil {
			return nil, fmt.Errorf("Docker daemon no disponible: %w", err)
		}
		a.dockerService = svc
	}
	return a.dockerService, nil
}

// GetOverview returns summary metrics of the Docker daemon
func (a *App) GetOverview() (*docker.SystemOverview, error) {
	svc, err := a.getService()
	if err != nil {
		return nil, err
	}
	return svc.GetOverview(a.ctx)
}

// ListContainers returns list of all containers
func (a *App) ListContainers(all bool) ([]docker.ContainerInfo, error) {
	svc, err := a.getService()
	if err != nil {
		return nil, err
	}
	return svc.ListContainers(a.ctx, all)
}

// StartContainer starts a container
func (a *App) StartContainer(id string) error {
	svc, err := a.getService()
	if err != nil {
		return err
	}
	return svc.StartContainer(a.ctx, id)
}

// StopContainer stops a container
func (a *App) StopContainer(id string) error {
	svc, err := a.getService()
	if err != nil {
		return err
	}
	return svc.StopContainer(a.ctx, id)
}

// RestartContainer restarts a container
func (a *App) RestartContainer(id string) error {
	svc, err := a.getService()
	if err != nil {
		return err
	}
	return svc.RestartContainer(a.ctx, id)
}

// StartStack starts all non-running containers in a compose project
func (a *App) StartStack(projectName string) error {
	svc, err := a.getService()
	if err != nil {
		return err
	}
	return svc.StartStack(a.ctx, projectName)
}

// StopStack stops all running containers in a compose project
func (a *App) StopStack(projectName string) error {
	svc, err := a.getService()
	if err != nil {
		return err
	}
	return svc.StopStack(a.ctx, projectName)
}

// RestartStack restarts all containers in a compose project
func (a *App) RestartStack(projectName string) error {
	svc, err := a.getService()
	if err != nil {
		return err
	}
	return svc.RestartStack(a.ctx, projectName)
}

// StartNetwork starts all non-running containers connected to a Docker network
func (a *App) StartNetwork(networkName string) error {
	svc, err := a.getService()
	if err != nil {
		return err
	}
	return svc.StartNetwork(a.ctx, networkName)
}

// StopNetwork stops all running containers connected to a Docker network
func (a *App) StopNetwork(networkName string) error {
	svc, err := a.getService()
	if err != nil {
		return err
	}
	return svc.StopNetwork(a.ctx, networkName)
}

// RestartNetwork restarts all containers connected to a Docker network
func (a *App) RestartNetwork(networkName string) error {
	svc, err := a.getService()
	if err != nil {
		return err
	}
	return svc.RestartNetwork(a.ctx, networkName)
}

// PauseContainer pauses a container
func (a *App) PauseContainer(id string) error {
	svc, err := a.getService()
	if err != nil {
		return err
	}
	return svc.PauseContainer(a.ctx, id)
}

// UnpauseContainer resumes a paused container
func (a *App) UnpauseContainer(id string) error {
	svc, err := a.getService()
	if err != nil {
		return err
	}
	return svc.UnpauseContainer(a.ctx, id)
}

// RemoveContainer deletes a container
func (a *App) RemoveContainer(id string, force bool) error {
	svc, err := a.getService()
	if err != nil {
		return err
	}
	return svc.RemoveContainer(a.ctx, id, force)
}

// GetContainerLogs returns logs for a container
func (a *App) GetContainerLogs(id string, tail int) (string, error) {
	svc, err := a.getService()
	if err != nil {
		return "", err
	}
	return svc.GetContainerLogs(a.ctx, id, tail)
}

// GetContainerStats returns current resource utilization
func (a *App) GetContainerStats(id string) (*docker.ContainerStats, error) {
	svc, err := a.getService()
	if err != nil {
		return nil, err
	}
	return svc.GetContainerStats(a.ctx, id)
}

// StartTerminal starts an interactive PTY shell inside a container
func (a *App) StartTerminal(containerID string, shell string, rows uint, cols uint) (*docker.TerminalStartResult, error) {
	svc, err := a.getService()
	if err != nil {
		return nil, err
	}

	return svc.StartTerminal(
		a.ctx,
		containerID,
		shell,
		rows,
		cols,
		func(sessID string, chunkBase64 string) {
			runtime.EventsEmit(a.ctx, "terminal:data:"+sessID, chunkBase64)
		},
		func(sessID string) {
			runtime.EventsEmit(a.ctx, "terminal:exit:"+sessID)
		},
	)
}

// WriteTerminal sends input data to an active container terminal session
func (a *App) WriteTerminal(sessionID string, data string) error {
	svc, err := a.getService()
	if err != nil {
		return err
	}
	return svc.WriteTerminal(sessionID, data)
}

// ResizeTerminal changes the window dimensions of a running container terminal
func (a *App) ResizeTerminal(sessionID string, rows uint, cols uint) error {
	svc, err := a.getService()
	if err != nil {
		return err
	}
	return svc.ResizeTerminal(a.ctx, sessionID, rows, cols)
}

// CloseTerminal terminates an active terminal session
func (a *App) CloseTerminal(sessionID string) error {
	a.mu.Lock()
	svc := a.dockerService
	a.mu.Unlock()

	if svc != nil {
		return svc.CloseTerminal(sessionID)
	}
	return nil
}

// ListImages returns list of all local images
func (a *App) ListImages() ([]docker.ImageInfo, error) {
	svc, err := a.getService()
	if err != nil {
		return nil, err
	}
	return svc.ListImages(a.ctx)
}

// GetDiskUsage returns aggregated metrics of image storage and reclaimable space
func (a *App) GetDiskUsage() (*docker.DiskUsageSummary, error) {
	svc, err := a.getService()
	if err != nil {
		return nil, err
	}
	return svc.GetDiskUsage(a.ctx)
}

// PullImage pulls an image and streams progress events
func (a *App) PullImage(imageName string) error {
	svc, err := a.getService()
	if err != nil {
		return err
	}

	return svc.PullImage(a.ctx, imageName, func(event docker.PullProgressEvent) {
		runtime.EventsEmit(a.ctx, "image:pull:progress", event)
	})
}

// RemoveImage removes an image by ID or name
func (a *App) RemoveImage(id string, force bool) error {
	svc, err := a.getService()
	if err != nil {
		return err
	}
	return svc.RemoveImage(a.ctx, id, force)
}

// PruneImages cleans dangling or unused images and returns reclaimed space
func (a *App) PruneImages(danglingOnly bool) (*docker.PruneResult, error) {
	svc, err := a.getService()
	if err != nil {
		return nil, err
	}
	return svc.PruneImages(a.ctx, danglingOnly)
}

// CreateContainer creates and optionally starts a new container
func (a *App) CreateContainer(req docker.CreateContainerRequest) (*docker.CreateContainerResult, error) {
	svc, err := a.getService()
	if err != nil {
		return nil, err
	}
	return svc.CreateContainer(a.ctx, req)
}
