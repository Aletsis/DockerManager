package network

import (
	"context"

	networkdomain "dockermanager/internal/domain/network"
)

// CreateNetworkUseCase handles the creation of Docker networks
type CreateNetworkUseCase struct {
	repo networkdomain.Repository
}

// NewCreateNetworkUseCase creates a new CreateNetworkUseCase
func NewCreateNetworkUseCase(repo networkdomain.Repository) *CreateNetworkUseCase {
	return &CreateNetworkUseCase{repo: repo}
}

// Execute validates and creates a new Docker network
func (uc *CreateNetworkUseCase) Execute(ctx context.Context, spec networkdomain.CreateNetworkSpec) (string, error) {
	if err := spec.Validate(); err != nil {
		return "", err
	}
	return uc.repo.Create(ctx, spec)
}
