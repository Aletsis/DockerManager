package stack

import (
	"context"

	stackdomain "dockermanager/internal/domain/stack"
)

// ManageStackUseCase handles compose project lifecycle mutations
type ManageStackUseCase struct {
	repo stackdomain.Repository
}

// NewManageStackUseCase creates a new ManageStackUseCase
func NewManageStackUseCase(repo stackdomain.Repository) *ManageStackUseCase {
	return &ManageStackUseCase{repo: repo}
}

// StartStack starts all stopped containers in a compose stack
func (uc *ManageStackUseCase) StartStack(ctx context.Context, projectName string) error {
	return uc.repo.StartStack(ctx, projectName)
}

// StopStack stops all running containers in a compose stack
func (uc *ManageStackUseCase) StopStack(ctx context.Context, projectName string) error {
	return uc.repo.StopStack(ctx, projectName)
}

// RestartStack restarts all containers in a compose stack
func (uc *ManageStackUseCase) RestartStack(ctx context.Context, projectName string) error {
	return uc.repo.RestartStack(ctx, projectName)
}
