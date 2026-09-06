package container

import (
	"context"

	containerdomain "dockermanager/internal/domain/container"
)

// CreateContainerUseCase validates and creates a new container
type CreateContainerUseCase struct {
	repo containerdomain.Repository
}

// NewCreateContainerUseCase creates a new CreateContainerUseCase
func NewCreateContainerUseCase(repo containerdomain.Repository) *CreateContainerUseCase {
	return &CreateContainerUseCase{repo: repo}
}

// Execute validates domain invariants and creates the container
func (uc *CreateContainerUseCase) Execute(ctx context.Context, spec containerdomain.CreateSpec) (*containerdomain.CreateResult, error) {
	if err := containerdomain.ValidateCreateSpec(spec); err != nil {
		return nil, err
	}
	return uc.repo.Create(ctx, spec)
}
