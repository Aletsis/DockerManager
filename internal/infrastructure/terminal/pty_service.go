package terminal

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	terminaldomain "dockermanager/internal/domain/terminal"
	"dockermanager/internal/infrastructure/docker"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Session represents an active terminal connection
type Session struct {
	ID          string
	ContainerID string
	ExecID      string
	Conn        net.Conn
	Cancel      context.CancelFunc
	closed      bool
	closeLock   sync.Mutex
}

// Close gracefully terminates the session
func (s *Session) Close() {
	s.closeLock.Lock()
	defer s.closeLock.Unlock()
	if s.closed {
		return
	}
	s.closed = true
	if s.Cancel != nil {
		s.Cancel()
	}
	if s.Conn != nil {
		_ = s.Conn.Close()
	}
}

// Service implements terminaldomain.Service using Docker Exec PTY
type Service struct {
	dockerClient *docker.Client
	mu           sync.RWMutex
	sessions     map[string]*Session
}

// NewService creates a new Terminal Service
func NewService(dockerClient *docker.Client) *Service {
	return &Service{
		dockerClient: dockerClient,
		sessions:     make(map[string]*Session),
	}
}

func generateSessionID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("term_%d_%s", time.Now().Unix(), hex.EncodeToString(b))
}

// resolveShell determines available shell binary in container
func resolveShell(ctx context.Context, cli *client.Client, containerID string, requestedShell string) string {
	clean := strings.TrimSpace(requestedShell)
	if clean != "" && clean != "auto" {
		return clean
	}

	if stat, err := cli.ContainerStatPath(ctx, containerID, "/bin/bash"); err == nil && !stat.Mode.IsDir() {
		return "/bin/bash"
	}
	if stat, err := cli.ContainerStatPath(ctx, containerID, "/bin/sh"); err == nil && !stat.Mode.IsDir() {
		return "/bin/sh"
	}
	return "/bin/sh"
}

// Start initiates a Docker exec process with PTY attached
func (s *Service) Start(
	ctx context.Context,
	containerID string,
	shell string,
	rows uint,
	cols uint,
	onOutput terminaldomain.OutputCallback,
	onExit terminaldomain.ExitCallback,
) (*terminaldomain.StartResult, error) {
	cli := s.dockerClient.RawClient()
	if cli == nil {
		return nil, errors.New("docker client not available")
	}

	if rows == 0 {
		rows = 24
	}
	if cols == 0 {
		cols = 80
	}

	resolvedShell := resolveShell(ctx, cli, containerID, shell)

	execConfig := container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{resolvedShell},
		ConsoleSize:  &[2]uint{rows, cols},
		Env: []string{
			"TERM=xterm-256color",
			"COLORTERM=truecolor",
		},
	}

	execResp, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create exec process: %w", err)
	}

	attachResp, err := cli.ContainerExecAttach(ctx, execResp.ID, container.ExecAttachOptions{
		Tty: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to attach to exec process: %w", err)
	}

	sessionCtx, cancel := context.WithCancel(context.Background())
	sessionID := generateSessionID()

	session := &Session{
		ID:          sessionID,
		ContainerID: containerID,
		ExecID:      execResp.ID,
		Conn:        attachResp.Conn,
		Cancel:      cancel,
	}

	s.mu.Lock()
	s.sessions[sessionID] = session
	s.mu.Unlock()

	_ = cli.ContainerExecResize(sessionCtx, execResp.ID, container.ResizeOptions{
		Height: rows,
		Width:  cols,
	})

	go func() {
		defer func() {
			_ = s.Close(sessionID)
			if onExit != nil {
				onExit(sessionID)
			}
		}()

		buf := make([]byte, 4096)
		for {
			select {
			case <-sessionCtx.Done():
				return
			default:
			}

			n, readErr := attachResp.Reader.Read(buf)
			if n > 0 && onOutput != nil {
				encoded := base64.StdEncoding.EncodeToString(buf[:n])
				onOutput(sessionID, encoded)
			}
			if readErr != nil {
				return
			}
		}
	}()

	return &terminaldomain.StartResult{
		SessionID: sessionID,
		Rows:      rows,
		Cols:      cols,
	}, nil
}

// Write writes data to active terminal PTY
func (s *Service) Write(sessionID string, data string) error {
	s.mu.RLock()
	session, exists := s.sessions[sessionID]
	s.mu.RUnlock()

	if !exists || session == nil {
		return errors.New("session not found")
	}

	_, err := session.Conn.Write([]byte(data))
	return err
}

// Resize updates terminal dimensions
func (s *Service) Resize(ctx context.Context, sessionID string, rows uint, cols uint) error {
	s.mu.RLock()
	session, exists := s.sessions[sessionID]
	s.mu.RUnlock()

	if !exists || session == nil {
		return errors.New("session not found")
	}

	if rows == 0 || cols == 0 {
		return nil
	}

	return s.dockerClient.RawClient().ContainerExecResize(ctx, session.ExecID, container.ResizeOptions{
		Height: rows,
		Width:  cols,
	})
}

// Close terminates a specific terminal session
func (s *Service) Close(sessionID string) error {
	s.mu.Lock()
	session, exists := s.sessions[sessionID]
	if exists {
		delete(s.sessions, sessionID)
	}
	s.mu.Unlock()

	if exists && session != nil {
		session.Close()
	}
	return nil
}

// CloseAll terminates all active sessions
func (s *Service) CloseAll() {
	s.mu.Lock()
	all := make([]*Session, 0, len(s.sessions))
	for _, sess := range s.sessions {
		all = append(all, sess)
	}
	s.sessions = make(map[string]*Session)
	s.mu.Unlock()

	for _, sess := range all {
		sess.Close()
	}
}
