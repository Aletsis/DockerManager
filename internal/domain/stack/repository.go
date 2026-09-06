package stack

import "context"

// Repository defines the port for compose stack lifecycle operations
type Repository interface {
	StartStack(ctx context.Context, projectName string) error
	StopStack(ctx context.Context, projectName string) error
	RestartStack(ctx context.Context, projectName string) error
}
