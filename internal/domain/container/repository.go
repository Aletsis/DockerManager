package container

import "context"

// Repository defines the port (contract) for container persistence and lifecycle operations
type Repository interface {
	List(ctx context.Context, all bool) ([]Container, error)
	FindByID(ctx context.Context, id string) (*Container, error)
	Start(ctx context.Context, id string) error
	Stop(ctx context.Context, id string) error
	Restart(ctx context.Context, id string) error
	Pause(ctx context.Context, id string) error
	Unpause(ctx context.Context, id string) error
	Remove(ctx context.Context, id string, force bool) error
	GetLogs(ctx context.Context, id string, tail int) (string, error)
	GetStats(ctx context.Context, id string) (*Stats, error)
	Create(ctx context.Context, spec CreateSpec) (*CreateResult, error)
}
