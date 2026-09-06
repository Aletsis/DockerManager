package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/image"
)

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
