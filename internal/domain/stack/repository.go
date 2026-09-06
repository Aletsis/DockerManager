package stack

import "context"

// Repository defines the port for compose stack lifecycle operations
type Repository interface {
	StartStack(ctx context.Context, projectName string) error
	StopStack(ctx context.Context, projectName string) error
	RestartStack(ctx context.Context, projectName string) error

	// Native Compose Operations
	UpStack(ctx context.Context, req ComposeDeployRequest, onLog func(string)) error
	DownStack(ctx context.Context, req ComposeDownRequest, onLog func(string)) error
	GetComposeFile(ctx context.Context, projectName, workingDir, configFile string) (*ComposeFileInfo, error)
	SaveComposeFile(ctx context.Context, workingDir, configFile, content string) (string, error)
	GetDefaultStackDirectory(projectName string) (string, error)
}

