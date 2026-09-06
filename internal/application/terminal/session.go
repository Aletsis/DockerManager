package terminal

import (
	"context"

	terminaldomain "dockermanager/internal/domain/terminal"
)

// ManageTerminalUseCase coordinates PTY terminal sessions
type ManageTerminalUseCase struct {
	service terminaldomain.Service
}

func NewManageTerminalUseCase(service terminaldomain.Service) *ManageTerminalUseCase {
	return &ManageTerminalUseCase{service: service}
}

func (uc *ManageTerminalUseCase) Start(
	ctx context.Context,
	containerID string,
	shell string,
	rows uint,
	cols uint,
	onOutput terminaldomain.OutputCallback,
	onExit terminaldomain.ExitCallback,
) (*terminaldomain.StartResult, error) {
	return uc.service.Start(ctx, containerID, shell, rows, cols, onOutput, onExit)
}

func (uc *ManageTerminalUseCase) Write(sessionID string, data string) error {
	return uc.service.Write(sessionID, data)
}

func (uc *ManageTerminalUseCase) Resize(ctx context.Context, sessionID string, rows uint, cols uint) error {
	return uc.service.Resize(ctx, sessionID, rows, cols)
}

func (uc *ManageTerminalUseCase) Close(sessionID string) error {
	return uc.service.Close(sessionID)
}
