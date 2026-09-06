package terminal

import "context"

// OutputCallback receives terminal output chunks (base64 encoded)
type OutputCallback func(sessionID string, chunkBase64 string)

// ExitCallback is invoked when the shell session terminates
type ExitCallback func(sessionID string)

// Service defines the port for terminal lifecycle operations
type Service interface {
	Start(ctx context.Context, containerID string, shell string, rows uint, cols uint, onOutput OutputCallback, onExit ExitCallback) (*StartResult, error)
	Write(sessionID string, data string) error
	Resize(ctx context.Context, sessionID string, rows uint, cols uint) error
	Close(sessionID string) error
}
