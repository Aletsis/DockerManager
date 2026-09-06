package system

import "context"

// Repository defines the port for system-level queries
type Repository interface {
	GetOverview(ctx context.Context) (*Overview, error)
}
