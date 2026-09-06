package main

import (
	"context"
	"fmt"
	"dockermanager/backend/docker"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct manages application state and services
type App struct {
	ctx           context.Context
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
	if a.dockerService != nil {
		_ = a.dockerService.Close()
	}
}

func (a *App) checkService() error {
	if a.dockerService == nil {
		svc, err := docker.NewService()
		if err != nil {
			return fmt.Errorf("Docker daemon no disponible: %w", err)
		}
		a.dockerService = svc
	}
	return nil
}

// GetOverview returns summary metrics of the Docker daemon
func (a *App) GetOverview() (*docker.SystemOverview, error) {
	if err := a.checkService(); err != nil {
		return nil, err
	}
	return a.dockerService.GetOverview(a.ctx)
}

// ListContainers returns list of all containers
func (a *App) ListContainers(all bool) ([]docker.ContainerInfo, error) {
	if err := a.checkService(); err != nil {
		return nil, err
	}
	return a.dockerService.ListContainers(a.ctx, all)
}

// StartContainer starts a container
func (a *App) StartContainer(id string) error {
	if err := a.checkService(); err != nil {
		return err
	}
	return a.dockerService.StartContainer(a.ctx, id)
}

// StopContainer stops a container
func (a *App) StopContainer(id string) error {
	if err := a.checkService(); err != nil {
		return err
	}
	return a.dockerService.StopContainer(a.ctx, id)
}

// RestartContainer restarts a container
func (a *App) RestartContainer(id string) error {
	if err := a.checkService(); err != nil {
		return err
	}
	return a.dockerService.RestartContainer(a.ctx, id)
}

// StartStack starts all non-running containers in a compose project
func (a *App) StartStack(projectName string) error {
	if err := a.checkService(); err != nil {
		return err
	}
	return a.dockerService.StartStack(a.ctx, projectName)
}

// StopStack stops all running containers in a compose project
func (a *App) StopStack(projectName string) error {
	if err := a.checkService(); err != nil {
		return err
	}
	return a.dockerService.StopStack(a.ctx, projectName)
}

// RestartStack restarts all containers in a compose project
func (a *App) RestartStack(projectName string) error {
	if err := a.checkService(); err != nil {
		return err
	}
	return a.dockerService.RestartStack(a.ctx, projectName)
}

// PauseContainer pauses a container
func (a *App) PauseContainer(id string) error {
	if err := a.checkService(); err != nil {
		return err
	}
	return a.dockerService.PauseContainer(a.ctx, id)
}

// UnpauseContainer resumes a paused container
func (a *App) UnpauseContainer(id string) error {
	if err := a.checkService(); err != nil {
		return err
	}
	return a.dockerService.UnpauseContainer(a.ctx, id)
}

// RemoveContainer deletes a container
func (a *App) RemoveContainer(id string, force bool) error {
	if err := a.checkService(); err != nil {
		return err
	}
	return a.dockerService.RemoveContainer(a.ctx, id, force)
}

// GetContainerLogs returns logs for a container
func (a *App) GetContainerLogs(id string, tail int) (string, error) {
	if err := a.checkService(); err != nil {
		return "", err
	}
	return a.dockerService.GetContainerLogs(a.ctx, id, tail)
}

// GetContainerStats returns current resource utilization
func (a *App) GetContainerStats(id string) (*docker.ContainerStats, error) {
	if err := a.checkService(); err != nil {
		return nil, err
	}
	return a.dockerService.GetContainerStats(a.ctx, id)
}

// StartTerminal starts an interactive PTY shell inside a container
func (a *App) StartTerminal(containerID string, shell string, rows uint, cols uint) (*docker.TerminalStartResult, error) {
	if err := a.checkService(); err != nil {
		return nil, err
	}

	return a.dockerService.StartTerminal(
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
	if err := a.checkService(); err != nil {
		return err
	}
	return a.dockerService.WriteTerminal(sessionID, data)
}

// ResizeTerminal changes the window dimensions of a running container terminal
func (a *App) ResizeTerminal(sessionID string, rows uint, cols uint) error {
	if err := a.checkService(); err != nil {
		return err
	}
	return a.dockerService.ResizeTerminal(a.ctx, sessionID, rows, cols)
}

// CloseTerminal terminates an active terminal session
func (a *App) CloseTerminal(sessionID string) error {
	if a.dockerService != nil {
		return a.dockerService.CloseTerminal(sessionID)
	}
	return nil
}

// ListImages returns list of all local images
func (a *App) ListImages() ([]docker.ImageInfo, error) {
	if err := a.checkService(); err != nil {
		return nil, err
	}
	return a.dockerService.ListImages(a.ctx)
}

// GetDiskUsage returns aggregated metrics of image storage and reclaimable space
func (a *App) GetDiskUsage() (*docker.DiskUsageSummary, error) {
	if err := a.checkService(); err != nil {
		return nil, err
	}
	return a.dockerService.GetDiskUsage(a.ctx)
}

// PullImage pulls an image and streams progress events
func (a *App) PullImage(imageName string) error {
	if err := a.checkService(); err != nil {
		return err
	}

	return a.dockerService.PullImage(a.ctx, imageName, func(event docker.PullProgressEvent) {
		runtime.EventsEmit(a.ctx, "image:pull:progress", event)
	})
}

// RemoveImage removes an image by ID or name
func (a *App) RemoveImage(id string, force bool) error {
	if err := a.checkService(); err != nil {
		return err
	}
	return a.dockerService.RemoveImage(a.ctx, id, force)
}

// PruneImages cleans dangling or unused images and returns reclaimed space
func (a *App) PruneImages(danglingOnly bool) (*docker.PruneResult, error) {
	if err := a.checkService(); err != nil {
		return nil, err
	}
	return a.dockerService.PruneImages(a.ctx, danglingOnly)
}

// CreateContainer creates and optionally starts a new container
func (a *App) CreateContainer(req docker.CreateContainerRequest) (*docker.CreateContainerResult, error) {
	if err := a.checkService(); err != nil {
		return nil, err
	}
	return a.dockerService.CreateContainer(a.ctx, req)
}


