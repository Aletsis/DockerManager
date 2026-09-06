package image

import (
	"context"

	imagedomain "dockermanager/internal/domain/image"
)

// ListImagesUseCase lists local images
type ListImagesUseCase struct {
	repo imagedomain.Repository
}

func NewListImagesUseCase(repo imagedomain.Repository) *ListImagesUseCase {
	return &ListImagesUseCase{repo: repo}
}

func (uc *ListImagesUseCase) Execute(ctx context.Context) ([]imagedomain.Image, error) {
	return uc.repo.List(ctx)
}

// GetDiskUsageUseCase retrieves storage metrics
type GetDiskUsageUseCase struct {
	repo imagedomain.Repository
}

func NewGetDiskUsageUseCase(repo imagedomain.Repository) *GetDiskUsageUseCase {
	return &GetDiskUsageUseCase{repo: repo}
}

func (uc *GetDiskUsageUseCase) Execute(ctx context.Context) (*imagedomain.DiskUsageSummary, error) {
	return uc.repo.GetDiskUsage(ctx)
}

// PullImageUseCase orchestrates image download
type PullImageUseCase struct {
	repo imagedomain.Repository
}

func NewPullImageUseCase(repo imagedomain.Repository) *PullImageUseCase {
	return &PullImageUseCase{repo: repo}
}

func (uc *PullImageUseCase) Execute(ctx context.Context, imageName string, onProgress func(imagedomain.PullProgressEvent)) error {
	return uc.repo.Pull(ctx, imageName, onProgress)
}

// RemoveImageUseCase removes an image
type RemoveImageUseCase struct {
	repo imagedomain.Repository
}

func NewRemoveImageUseCase(repo imagedomain.Repository) *RemoveImageUseCase {
	return &RemoveImageUseCase{repo: repo}
}

func (uc *RemoveImageUseCase) Execute(ctx context.Context, id string, force bool) error {
	return uc.repo.Remove(ctx, id, force)
}

// PruneImagesUseCase cleans unused images
type PruneImagesUseCase struct {
	repo imagedomain.Repository
}

func NewPruneImagesUseCase(repo imagedomain.Repository) *PruneImagesUseCase {
	return &PruneImagesUseCase{repo: repo}
}

func (uc *PruneImagesUseCase) Execute(ctx context.Context, danglingOnly bool) (*imagedomain.PruneResult, error) {
	return uc.repo.Prune(ctx, danglingOnly)
}
