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
