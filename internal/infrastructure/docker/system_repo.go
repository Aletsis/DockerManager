package docker

import (
	"context"
	"fmt"

	systemdomain "dockermanager/internal/domain/system"

	"github.com/docker/docker/api/types/image"
)

// SystemRepository implements systemdomain.Repository using Docker SDK
type SystemRepository struct {
	client *Client
}

// NewSystemRepository creates a new Docker SystemRepository
func NewSystemRepository(client *Client) *SystemRepository {
	return &SystemRepository{client: client}
}

// GetOverview returns system and Docker daemon metrics
func (r *SystemRepository) GetOverview(ctx context.Context) (*systemdomain.Overview, error) {
	info, err := r.client.cli.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve docker info: %w", err)
	}

	version, err := r.client.cli.ServerVersion(ctx)
	versionStr := ""
	if err == nil {
		versionStr = version.Version
	}

	images, _ := r.client.cli.ImageList(ctx, image.ListOptions{})

	return &systemdomain.Overview{
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
