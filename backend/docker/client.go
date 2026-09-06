package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
)

// Service encapsulates Docker client operations and coordinates domain services
type Service struct {
	cli             *client.Client
	terminalManager *TerminalManager
}

// NewService instantiates and connects a Docker client
func NewService() (*Service, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}
	return &Service{
		cli:             cli,
		terminalManager: NewTerminalManager(),
	}, nil
}

// Close closes the Docker client connection and cleans up sessions
func (s *Service) Close() error {
	if s.terminalManager != nil {
		s.terminalManager.CloseAll()
	}
	if s.cli != nil {
		return s.cli.Close()
	}
	return nil
}

// StartTerminal starts an interactive PTY session for a container
func (s *Service) StartTerminal(
	ctx context.Context,
	containerID string,
	shell string,
	rows uint,
	cols uint,
	onData func(sessionID string, chunkBase64 string),
	onExit func(sessionID string),
) (*TerminalStartResult, error) {
	if s.terminalManager == nil {
		return nil, fmt.Errorf("terminal manager not initialized")
	}
	return s.terminalManager.StartSession(ctx, s.cli, containerID, shell, rows, cols, onData, onExit)
}

// WriteTerminal sends input data to a terminal session
func (s *Service) WriteTerminal(sessionID string, data string) error {
	if s.terminalManager == nil {
		return fmt.Errorf("terminal manager not initialized")
	}
	return s.terminalManager.Write(sessionID, data)
}

// ResizeTerminal updates the terminal window dimensions
func (s *Service) ResizeTerminal(ctx context.Context, sessionID string, rows uint, cols uint) error {
	if s.terminalManager == nil {
		return fmt.Errorf("terminal manager not initialized")
	}
	return s.terminalManager.Resize(ctx, s.cli, sessionID, rows, cols)
}

// CloseTerminal terminates an active terminal session
func (s *Service) CloseTerminal(sessionID string) error {
	if s.terminalManager == nil {
		return nil
	}
	s.terminalManager.CloseSession(sessionID)
	return nil
}
