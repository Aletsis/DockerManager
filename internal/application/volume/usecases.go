package volume

import (
	"context"

	volumedomain "dockermanager/internal/domain/volume"
)

// ListVolumesUseCase lists all volumes with usage information
type ListVolumesUseCase struct {
	repo volumedomain.Repository
}

func NewListVolumesUseCase(repo volumedomain.Repository) *ListVolumesUseCase {
	return &ListVolumesUseCase{repo: repo}
}

func (uc *ListVolumesUseCase) Execute(ctx context.Context) ([]volumedomain.Volume, error) {
	return uc.repo.List(ctx)
}

// GetVolumeDiskUsageUseCase returns disk metrics for volumes
type GetVolumeDiskUsageUseCase struct {
	repo volumedomain.Repository
}

func NewGetVolumeDiskUsageUseCase(repo volumedomain.Repository) *GetVolumeDiskUsageUseCase {
	return &GetVolumeDiskUsageUseCase{repo: repo}
}

func (uc *GetVolumeDiskUsageUseCase) Execute(ctx context.Context) (*volumedomain.DiskUsageSummary, error) {
	return uc.repo.GetDiskUsage(ctx)
}

// InspectVolumeUseCase inspects a volume returning raw formatted JSON
type InspectVolumeUseCase struct {
	repo volumedomain.Repository
}

func NewInspectVolumeUseCase(repo volumedomain.Repository) *InspectVolumeUseCase {
	return &InspectVolumeUseCase{repo: repo}
}

func (uc *InspectVolumeUseCase) Execute(ctx context.Context, name string) (string, error) {
	return uc.repo.Inspect(ctx, name)
}

// RemoveVolumeUseCase removes a volume
type RemoveVolumeUseCase struct {
	repo volumedomain.Repository
}

func NewRemoveVolumeUseCase(repo volumedomain.Repository) *RemoveVolumeUseCase {
	return &RemoveVolumeUseCase{repo: repo}
}

func (uc *RemoveVolumeUseCase) Execute(ctx context.Context, name string, force bool) error {
	return uc.repo.Remove(ctx, name, force)
}

// PruneVolumesUseCase cleans unused/dangling volumes
type PruneVolumesUseCase struct {
	repo volumedomain.Repository
}

func NewPruneVolumesUseCase(repo volumedomain.Repository) *PruneVolumesUseCase {
	return &PruneVolumesUseCase{repo: repo}
}

func (uc *PruneVolumesUseCase) Execute(ctx context.Context) (*volumedomain.PruneResult, error) {
	return uc.repo.Prune(ctx)
}
