package network

import (
	"context"

	networkdomain "dockermanager/internal/domain/network"
)

// ListNetworksUseCase handles retrieving all networks
type ListNetworksUseCase struct {
	repo networkdomain.Repository
}

// NewListNetworksUseCase creates a new ListNetworksUseCase
func NewListNetworksUseCase(repo networkdomain.Repository) *ListNetworksUseCase {
	return &ListNetworksUseCase{repo: repo}
}

// Execute retrieves all docker networks
func (uc *ListNetworksUseCase) Execute(ctx context.Context) ([]networkdomain.Network, error) {
	return uc.repo.List(ctx)
}
