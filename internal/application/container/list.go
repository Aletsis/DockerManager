package container

import (
	"context"

	containerdomain "dockermanager/internal/domain/container"
)

// ListContainersUseCase retrieves all or active containers
type ListContainersUseCase struct {
	repo containerdomain.Repository
}

// NewListContainersUseCase creates a new ListContainersUseCase
func NewListContainersUseCase(repo containerdomain.Repository) *ListContainersUseCase {
	return &ListContainersUseCase{repo: repo}
}

// Execute runs the use case
func (uc *ListContainersUseCase) Execute(ctx context.Context, all bool) ([]containerdomain.Container, error) {
	return uc.repo.List(ctx, all)
}
