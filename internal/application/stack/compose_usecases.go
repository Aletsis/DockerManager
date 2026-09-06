package stack

import (
	"context"

	stackdomain "dockermanager/internal/domain/stack"
)

// UpStack deploys or updates a compose stack using native docker compose
func (uc *ManageStackUseCase) UpStack(ctx context.Context, req stackdomain.ComposeDeployRequest, onLog func(string)) error {
	return uc.repo.UpStack(ctx, req, onLog)
}

// DownStack tears down a compose stack
func (uc *ManageStackUseCase) DownStack(ctx context.Context, req stackdomain.ComposeDownRequest, onLog func(string)) error {
	return uc.repo.DownStack(ctx, req, onLog)
}

// GetComposeFile loads the compose file details
func (uc *ManageStackUseCase) GetComposeFile(ctx context.Context, projectName, workingDir, configFile string) (*stackdomain.ComposeFileInfo, error) {
	return uc.repo.GetComposeFile(ctx, projectName, workingDir, configFile)
}

// SaveComposeFile saves compose content to disk
func (uc *ManageStackUseCase) SaveComposeFile(ctx context.Context, workingDir, configFile, content string) (string, error) {
	return uc.repo.SaveComposeFile(ctx, workingDir, configFile, content)
}

// GetDefaultStackDirectory returns default path ~/.dockermanager/stacks/<projectName>
func (uc *ManageStackUseCase) GetDefaultStackDirectory(projectName string) (string, error) {
	return uc.repo.GetDefaultStackDirectory(projectName)
}
