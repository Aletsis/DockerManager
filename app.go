package main

import (
	"context"
	"fmt"
	"sync"

	containerapp "dockermanager/internal/application/container"
	imageapp "dockermanager/internal/application/image"
	networkapp "dockermanager/internal/application/network"
	stackapp "dockermanager/internal/application/stack"
	systemapp "dockermanager/internal/application/system"
	terminalapp "dockermanager/internal/application/terminal"
	volumeapp "dockermanager/internal/application/volume"

	containerdomain "dockermanager/internal/domain/container"
	imagedomain "dockermanager/internal/domain/image"
	networkdomain "dockermanager/internal/domain/network"
	systemdomain "dockermanager/internal/domain/system"
	terminaldomain "dockermanager/internal/domain/terminal"
	volumedomain "dockermanager/internal/domain/volume"

	dockerinfra "dockermanager/internal/infrastructure/docker"
	terminalinfra "dockermanager/internal/infrastructure/terminal"
	"dockermanager/internal/version"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App is the Wails delivery controller orchestrating application use cases
type App struct {
	ctx context.Context
	mu  sync.Mutex

	// Infrastructure Adapters
	dockerClient    *dockerinfra.Client
	terminalService *terminalinfra.Service

	// Application Use Cases
	listContainersUC *containerapp.ListContainersUseCase
	lifecycleUC      *containerapp.ManageLifecycleUseCase
	createUC         *containerapp.CreateContainerUseCase
	telemetryUC      *containerapp.GetTelemetryUseCase

	stackUC         *stackapp.ManageStackUseCase
	networkUC       *networkapp.ManageNetworkUseCase
	listNetworksUC  *networkapp.ListNetworksUseCase
	createNetworkUC *networkapp.CreateNetworkUseCase

	listImagesUC  *imageapp.ListImagesUseCase
	diskUsageUC   *imageapp.GetDiskUsageUseCase
	pullImageUC   *imageapp.PullImageUseCase
	removeImageUC *imageapp.RemoveImageUseCase
	pruneImagesUC *imageapp.PruneImagesUseCase

	listVolumesUC     *volumeapp.ListVolumesUseCase
	volumeDiskUsageUC *volumeapp.GetVolumeDiskUsageUseCase
	inspectVolumeUC   *volumeapp.InspectVolumeUseCase
	removeVolumeUC    *volumeapp.RemoveVolumeUseCase
	pruneVolumesUC    *volumeapp.PruneVolumesUseCase

	overviewUC *systemapp.GetOverviewUseCase
	terminalUC *terminalapp.ManageTerminalUseCase
}

// NewApp creates a new App controller
func NewApp() *App {
	return &App{}
}

// startup is called by Wails when the application boots
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	if err := a.ensureInitialized(); err != nil {
		fmt.Println("Warning: Failed to initialize Docker infrastructure:", err)
	}
}

// shutdown is called by Wails when the application terminates
func (a *App) shutdown(ctx context.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.terminalService != nil {
		a.terminalService.CloseAll()
	}
	if a.dockerClient != nil {
		_ = a.dockerClient.Close()
		a.dockerClient = nil
	}
}

// ensureInitialized lazily wires infrastructure and use cases with thread-safety
func (a *App) ensureInitialized() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.dockerClient != nil {
		return nil
	}

	cli, err := dockerinfra.NewClient()
	if err != nil {
		return fmt.Errorf("Docker daemon no disponible: %w", err)
	}
	a.dockerClient = cli

	// Domain Ports implemented by Infrastructure Adapters
	containerRepo := dockerinfra.NewContainerRepository(cli)
	imageRepo := dockerinfra.NewImageRepository(cli)
	volumeRepo := dockerinfra.NewVolumeRepository(cli)
	stackRepo := dockerinfra.NewStackRepository(cli)
	networkRepo := dockerinfra.NewNetworkRepository(cli)
	systemRepo := dockerinfra.NewSystemRepository(cli)
	a.terminalService = terminalinfra.NewService(cli)

	// Application Layer Use Cases
	a.listContainersUC = containerapp.NewListContainersUseCase(containerRepo)
	a.lifecycleUC = containerapp.NewManageLifecycleUseCase(containerRepo)
	a.createUC = containerapp.NewCreateContainerUseCase(containerRepo)
	a.telemetryUC = containerapp.NewGetTelemetryUseCase(containerRepo)

	a.stackUC = stackapp.NewManageStackUseCase(stackRepo)
	a.networkUC = networkapp.NewManageNetworkUseCase(networkRepo)
	a.listNetworksUC = networkapp.NewListNetworksUseCase(networkRepo)
	a.createNetworkUC = networkapp.NewCreateNetworkUseCase(networkRepo)

	a.listImagesUC = imageapp.NewListImagesUseCase(imageRepo)
	a.diskUsageUC = imageapp.NewGetDiskUsageUseCase(imageRepo)
	a.pullImageUC = imageapp.NewPullImageUseCase(imageRepo)
	a.removeImageUC = imageapp.NewRemoveImageUseCase(imageRepo)
	a.pruneImagesUC = imageapp.NewPruneImagesUseCase(imageRepo)

	a.listVolumesUC = volumeapp.NewListVolumesUseCase(volumeRepo)
	a.volumeDiskUsageUC = volumeapp.NewGetVolumeDiskUsageUseCase(volumeRepo)
	a.inspectVolumeUC = volumeapp.NewInspectVolumeUseCase(volumeRepo)
	a.removeVolumeUC = volumeapp.NewRemoveVolumeUseCase(volumeRepo)
	a.pruneVolumesUC = volumeapp.NewPruneVolumesUseCase(volumeRepo)

	a.overviewUC = systemapp.NewGetOverviewUseCase(systemRepo)
	a.terminalUC = terminalapp.NewManageTerminalUseCase(a.terminalService)

	return nil
}

// GetAppVersion returns the current application version of DockerManager
func (a *App) GetAppVersion() string {
	return version.Current
}

// GetOverview returns summary metrics of the Docker host
func (a *App) GetOverview() (*systemdomain.Overview, error) {
	if err := a.ensureInitialized(); err != nil {
		return nil, err
	}
	return a.overviewUC.Execute(a.ctx)
}

// ListContainers returns list of containers
func (a *App) ListContainers(all bool) ([]containerdomain.Container, error) {
	if err := a.ensureInitialized(); err != nil {
		return nil, err
	}
	return a.listContainersUC.Execute(a.ctx, all)
}

// StartContainer starts a container
func (a *App) StartContainer(id string) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.lifecycleUC.Start(a.ctx, id)
}

// StopContainer stops a container
func (a *App) StopContainer(id string) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.lifecycleUC.Stop(a.ctx, id)
}

// RestartContainer restarts a container
func (a *App) RestartContainer(id string) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.lifecycleUC.Restart(a.ctx, id)
}

// PauseContainer pauses a container
func (a *App) PauseContainer(id string) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.lifecycleUC.Pause(a.ctx, id)
}

// UnpauseContainer resumes a paused container
func (a *App) UnpauseContainer(id string) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.lifecycleUC.Unpause(a.ctx, id)
}

// RemoveContainer deletes a container
func (a *App) RemoveContainer(id string, force bool) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.lifecycleUC.Remove(a.ctx, id, force)
}

// StartStack starts all non-running containers in a compose stack
func (a *App) StartStack(projectName string) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.stackUC.StartStack(a.ctx, projectName)
}

// StopStack stops all running containers in a compose stack
func (a *App) StopStack(projectName string) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.stackUC.StopStack(a.ctx, projectName)
}

// RestartStack restarts all containers in a compose stack
func (a *App) RestartStack(projectName string) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.stackUC.RestartStack(a.ctx, projectName)
}

// StartNetwork starts all non-running containers attached to a network
func (a *App) StartNetwork(networkName string) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.networkUC.StartNetwork(a.ctx, networkName)
}

// StopNetwork stops all running containers attached to a network
func (a *App) StopNetwork(networkName string) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.networkUC.StopNetwork(a.ctx, networkName)
}

// RestartNetwork restarts all containers attached to a network
func (a *App) RestartNetwork(networkName string) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.networkUC.RestartNetwork(a.ctx, networkName)
}

// ListNetworks returns list of all docker networks
func (a *App) ListNetworks() ([]networkdomain.Network, error) {
	if err := a.ensureInitialized(); err != nil {
		return nil, err
	}
	return a.listNetworksUC.Execute(a.ctx)
}

// CreateNetwork creates a new docker network
func (a *App) CreateNetwork(spec networkdomain.CreateNetworkSpec) (string, error) {
	if err := a.ensureInitialized(); err != nil {
		return "", err
	}
	return a.createNetworkUC.Execute(a.ctx, spec)
}

// RemoveNetwork removes a docker network by id or name
func (a *App) RemoveNetwork(idOrName string) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.networkUC.Remove(a.ctx, idOrName)
}

// PruneNetworks cleans inactive/unused networks
func (a *App) PruneNetworks() (*networkdomain.PruneResult, error) {
	if err := a.ensureInitialized(); err != nil {
		return nil, err
	}
	return a.networkUC.Prune(a.ctx)
}

// ConnectContainerToNetwork attaches a container to a network
func (a *App) ConnectContainerToNetwork(networkID string, containerID string, ipAddress string) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.networkUC.Connect(a.ctx, networkID, containerID, ipAddress)
}

// DisconnectContainerFromNetwork detaches a container from a network
func (a *App) DisconnectContainerFromNetwork(networkID string, containerID string, force bool) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.networkUC.Disconnect(a.ctx, networkID, containerID, force)
}

// InspectNetwork returns raw formatted JSON configuration of a network
func (a *App) InspectNetwork(idOrName string) (string, error) {
	if err := a.ensureInitialized(); err != nil {
		return "", err
	}
	return a.networkUC.Inspect(a.ctx, idOrName)
}

// GetContainerLogs returns logs for a container
func (a *App) GetContainerLogs(id string, tail int) (string, error) {
	if err := a.ensureInitialized(); err != nil {
		return "", err
	}
	return a.telemetryUC.GetLogs(a.ctx, id, tail)
}

// GetContainerStats returns current resource utilization
func (a *App) GetContainerStats(id string) (*containerdomain.Stats, error) {
	if err := a.ensureInitialized(); err != nil {
		return nil, err
	}
	return a.telemetryUC.GetStats(a.ctx, id)
}

// InspectContainer returns raw formatted JSON configuration of a container
func (a *App) InspectContainer(id string) (string, error) {
	if err := a.ensureInitialized(); err != nil {
		return "", err
	}
	return a.telemetryUC.Inspect(a.ctx, id)
}

// OpenURL opens the given URL in the user's default web browser
func (a *App) OpenURL(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}

// StartTerminal starts an interactive PTY shell inside a container
func (a *App) StartTerminal(containerID string, shell string, rows uint, cols uint) (*terminaldomain.StartResult, error) {
	if err := a.ensureInitialized(); err != nil {
		return nil, err
	}

	return a.terminalUC.Start(
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
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.terminalUC.Write(sessionID, data)
}

// ResizeTerminal changes the window dimensions of a running container terminal
func (a *App) ResizeTerminal(sessionID string, rows uint, cols uint) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.terminalUC.Resize(a.ctx, sessionID, rows, cols)
}

// CloseTerminal terminates an active terminal session
func (a *App) CloseTerminal(sessionID string) error {
	a.mu.Lock()
	tuc := a.terminalUC
	a.mu.Unlock()

	if tuc != nil {
		return tuc.Close(sessionID)
	}
	return nil
}

// ListImages returns list of local images
func (a *App) ListImages() ([]imagedomain.Image, error) {
	if err := a.ensureInitialized(); err != nil {
		return nil, err
	}
	return a.listImagesUC.Execute(a.ctx)
}

// GetDiskUsage returns aggregated metrics of image storage
func (a *App) GetDiskUsage() (*imagedomain.DiskUsageSummary, error) {
	if err := a.ensureInitialized(); err != nil {
		return nil, err
	}
	return a.diskUsageUC.Execute(a.ctx)
}

// PullImage pulls an image and streams progress events
func (a *App) PullImage(imageName string) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}

	return a.pullImageUC.Execute(a.ctx, imageName, func(event imagedomain.PullProgressEvent) {
		runtime.EventsEmit(a.ctx, "image:pull:progress", event)
	})
}

// RemoveImage removes an image by ID or name
func (a *App) RemoveImage(id string, force bool) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.removeImageUC.Execute(a.ctx, id, force)
}

// PruneImages cleans dangling or unused images
func (a *App) PruneImages(danglingOnly bool) (*imagedomain.PruneResult, error) {
	if err := a.ensureInitialized(); err != nil {
		return nil, err
	}
	return a.pruneImagesUC.Execute(a.ctx, danglingOnly)
}

// CreateContainer creates and optionally starts a new container
func (a *App) CreateContainer(req containerdomain.CreateSpec) (*containerdomain.CreateResult, error) {
	if err := a.ensureInitialized(); err != nil {
		return nil, err
	}
	return a.createUC.Execute(a.ctx, req)
}

// ListVolumes returns list of all volumes with usage data
func (a *App) ListVolumes() ([]volumedomain.Volume, error) {
	if err := a.ensureInitialized(); err != nil {
		return nil, err
	}
	return a.listVolumesUC.Execute(a.ctx)
}

// GetVolumeDiskUsage returns aggregated metrics of volume storage
func (a *App) GetVolumeDiskUsage() (*volumedomain.DiskUsageSummary, error) {
	if err := a.ensureInitialized(); err != nil {
		return nil, err
	}
	return a.volumeDiskUsageUC.Execute(a.ctx)
}

// InspectVolume returns raw formatted JSON configuration of a volume
func (a *App) InspectVolume(name string) (string, error) {
	if err := a.ensureInitialized(); err != nil {
		return "", err
	}
	return a.inspectVolumeUC.Execute(a.ctx, name)
}

// RemoveVolume removes a volume by name
func (a *App) RemoveVolume(name string, force bool) error {
	if err := a.ensureInitialized(); err != nil {
		return err
	}
	return a.removeVolumeUC.Execute(a.ctx, name, force)
}

// PruneVolumes cleans dangling or unused volumes
func (a *App) PruneVolumes() (*volumedomain.PruneResult, error) {
	if err := a.ensureInitialized(); err != nil {
		return nil, err
	}
	return a.pruneVolumesUC.Execute(a.ctx)
}

