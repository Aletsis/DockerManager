package container

import (
	"context"

	containerdomain "dockermanager/internal/domain/container"
)

// ManageLifecycleUseCase handles container state mutations
type ManageLifecycleUseCase struct {
	repo containerdomain.Repository
}

// NewManageLifecycleUseCase creates a new ManageLifecycleUseCase
func NewManageLifecycleUseCase(repo containerdomain.Repository) *ManageLifecycleUseCase {
	return &ManageLifecycleUseCase{repo: repo}
}

// Start starts a container
func (uc *ManageLifecycleUseCase) Start(ctx context.Context, id string) error {
	return uc.repo.Start(ctx, id)
}

// Stop stops a container
func (uc *ManageLifecycleUseCase) Stop(ctx context.Context, id string) error {
	return uc.repo.Stop(ctx, id)
}

// Restart restarts a container
func (uc *ManageLifecycleUseCase) Restart(ctx context.Context, id string) error {
	return uc.repo.Restart(ctx, id)
}

// Pause pauses a container
func (uc *ManageLifecycleUseCase) Pause(ctx context.Context, id string) error {
	return uc.repo.Pause(ctx, id)
}

// Unpause unpauses a container
func (uc *ManageLifecycleUseCase) Unpause(ctx context.Context, id string) error {
	return uc.repo.Unpause(ctx, id)
}

// Remove deletes a container
func (uc *ManageLifecycleUseCase) Remove(ctx context.Context, id string, force bool) error {
	return uc.repo.Remove(ctx, id, force)
}
