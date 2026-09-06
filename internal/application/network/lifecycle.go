package network

import (
	"context"

	networkdomain "dockermanager/internal/domain/network"
)

// ManageNetworkUseCase handles network-wide container operations
type ManageNetworkUseCase struct {
	repo networkdomain.Repository
}

// NewManageNetworkUseCase creates a new ManageNetworkUseCase
func NewManageNetworkUseCase(repo networkdomain.Repository) *ManageNetworkUseCase {
	return &ManageNetworkUseCase{repo: repo}
}

// StartNetwork starts all containers attached to a network
func (uc *ManageNetworkUseCase) StartNetwork(ctx context.Context, networkName string) error {
	return uc.repo.StartNetwork(ctx, networkName)
}

// StopNetwork stops all containers attached to a network
func (uc *ManageNetworkUseCase) StopNetwork(ctx context.Context, networkName string) error {
	return uc.repo.StopNetwork(ctx, networkName)
}

// RestartNetwork restarts all containers attached to a network
func (uc *ManageNetworkUseCase) RestartNetwork(ctx context.Context, networkName string) error {
	return uc.repo.RestartNetwork(ctx, networkName)
}

// Remove deletes a network
func (uc *ManageNetworkUseCase) Remove(ctx context.Context, idOrName string) error {
	return uc.repo.Remove(ctx, idOrName)
}

// Prune cleans all unused networks
func (uc *ManageNetworkUseCase) Prune(ctx context.Context) (*networkdomain.PruneResult, error) {
	return uc.repo.Prune(ctx)
}

// Connect attaches a container to a network
func (uc *ManageNetworkUseCase) Connect(ctx context.Context, networkID, containerID, ipAddress string) error {
	return uc.repo.ConnectContainer(ctx, networkID, containerID, ipAddress)
}

// Disconnect detaches a container from a network
func (uc *ManageNetworkUseCase) Disconnect(ctx context.Context, networkID, containerID string, force bool) error {
	return uc.repo.DisconnectContainer(ctx, networkID, containerID, force)
}

// Inspect returns raw formatted JSON of a network
func (uc *ManageNetworkUseCase) Inspect(ctx context.Context, idOrName string) (string, error) {
	return uc.repo.Inspect(ctx, idOrName)
}
