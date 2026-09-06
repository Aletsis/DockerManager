package image

import "context"

// Repository defines the port for image operations
type Repository interface {
	List(ctx context.Context) ([]Image, error)
	GetDiskUsage(ctx context.Context) (*DiskUsageSummary, error)
	Pull(ctx context.Context, imageName string, onProgress func(PullProgressEvent)) error
	Remove(ctx context.Context, id string, force bool) error
	Prune(ctx context.Context, danglingOnly bool) (*PruneResult, error)
}
