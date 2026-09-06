package network

import "context"

// Repository defines the port for network-level operations
type Repository interface {
	StartNetwork(ctx context.Context, networkName string) error
	StopNetwork(ctx context.Context, networkName string) error
	RestartNetwork(ctx context.Context, networkName string) error
}
