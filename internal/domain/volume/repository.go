package volume

import "context"

// Repository defines the port interface for volume operations
type Repository interface {
	List(ctx context.Context) ([]Volume, error)
	GetDiskUsage(ctx context.Context) (*DiskUsageSummary, error)
	Inspect(ctx context.Context, name string) (string, error)
	Remove(ctx context.Context, name string, force bool) error
	Prune(ctx context.Context) (*PruneResult, error)
}
