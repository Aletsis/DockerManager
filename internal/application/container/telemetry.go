package container

import (
	"context"

	containerdomain "dockermanager/internal/domain/container"
)

// GetTelemetryUseCase handles logs and stats queries
type GetTelemetryUseCase struct {
	repo containerdomain.Repository
}

// NewGetTelemetryUseCase creates a new GetTelemetryUseCase
func NewGetTelemetryUseCase(repo containerdomain.Repository) *GetTelemetryUseCase {
	return &GetTelemetryUseCase{repo: repo}
}

// GetLogs returns container log stream
func (uc *GetTelemetryUseCase) GetLogs(ctx context.Context, id string, tail int) (string, error) {
	return uc.repo.GetLogs(ctx, id, tail)
}

// GetStats returns current container resource metrics
func (uc *GetTelemetryUseCase) GetStats(ctx context.Context, id string) (*containerdomain.Stats, error) {
	return uc.repo.GetStats(ctx, id)
}
