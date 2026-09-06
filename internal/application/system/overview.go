package system

import (
	"context"

	systemdomain "dockermanager/internal/domain/system"
)

// GetOverviewUseCase retrieves system-wide summary metrics
type GetOverviewUseCase struct {
	repo systemdomain.Repository
}

func NewGetOverviewUseCase(repo systemdomain.Repository) *GetOverviewUseCase {
	return &GetOverviewUseCase{repo: repo}
}

func (uc *GetOverviewUseCase) Execute(ctx context.Context) (*systemdomain.Overview, error) {
	return uc.repo.GetOverview(ctx)
}
