package network

import "context"

// Repository defines the port for network-level operations
type Repository interface {
	List(ctx context.Context) ([]Network, error)
	Inspect(ctx context.Context, idOrName string) (string, error)
	Create(ctx context.Context, spec CreateNetworkSpec) (string, error)
	Remove(ctx context.Context, idOrName string) error
	Prune(ctx context.Context) (*PruneResult, error)
	ConnectContainer(ctx context.Context, networkID, containerID string, ipAddress string) error
	DisconnectContainer(ctx context.Context, networkID, containerID string, force bool) error
	StartNetwork(ctx context.Context, networkName string) error
	StopNetwork(ctx context.Context, networkName string) error
	RestartNetwork(ctx context.Context, networkName string) error
}
