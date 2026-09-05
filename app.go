package main

import (
	"context"
	"fmt"
	"dockermanager/backend/docker"
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
