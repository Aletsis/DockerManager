package docker

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

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// TerminalStartResult contains information about a started terminal session
type TerminalStartResult struct {
	SessionID string `json:"sessionId"`
	Shell     string `json:"shell"`
}

// TerminalSession encapsulates an active PTY exec connection
type TerminalSession struct {
	ID          string
	ContainerID string
	ExecID      string
	Conn        net.Conn
	Cancel      context.CancelFunc
	closed      bool
	closeLock   sync.Mutex
}

// Close gracefully terminates the terminal session
func (s *TerminalSession) Close() {
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

// TerminalManager tracks all active terminal sessions
type TerminalManager struct {
	mu       sync.RWMutex
	sessions map[string]*TerminalSession
}

// NewTerminalManager creates a new TerminalManager
func NewTerminalManager() *TerminalManager {
	return &TerminalManager{
		sessions: make(map[string]*TerminalSession),
	}
}

func generateSessionID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("term_%d_%s", time.Now().Unix(), hex.EncodeToString(b))
}

// ResolveShell determines which shell binary to execute in the container
func ResolveShell(ctx context.Context, cli *client.Client, containerID string, requestedShell string) string {
	clean := strings.TrimSpace(requestedShell)
	if clean != "" && clean != "auto" {
		return clean
	}

	// Try /bin/bash first
	if stat, err := cli.ContainerStatPath(ctx, containerID, "/bin/bash"); err == nil && !stat.Mode.IsDir() {
		return "/bin/bash"
	}

	// Fallback to /bin/sh
	if stat, err := cli.ContainerStatPath(ctx, containerID, "/bin/sh"); err == nil && !stat.Mode.IsDir() {
		return "/bin/sh"
	}

	// Default fallback
	return "/bin/sh"
}

// StartSession initiates a Docker exec process with PTY attached
func (tm *TerminalManager) StartSession(
	ctx context.Context,
	cli *client.Client,
	containerID string,
	shell string,
	rows uint,
	cols uint,
	onData func(sessionID string, base64Chunk string),
	onExit func(sessionID string),
) (*TerminalStartResult, error) {
	if cli == nil {
		return nil, errors.New("docker client not available")
	}

	if rows == 0 {
		rows = 24
	}
	if cols == 0 {
		cols = 80
	}

	resolvedShell := ResolveShell(ctx, cli, containerID, shell)

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

	session := &TerminalSession{
		ID:          sessionID,
		ContainerID: containerID,
		ExecID:      execResp.ID,
		Conn:        attachResp.Conn,
		Cancel:      cancel,
	}

	tm.mu.Lock()
	tm.sessions[sessionID] = session
	tm.mu.Unlock()

	// Initial resize to ensure dimensions match
	_ = cli.ContainerExecResize(sessionCtx, execResp.ID, container.ResizeOptions{
		Height: rows,
		Width:  cols,
	})

	// Background reader goroutine
	go func() {
		defer func() {
			tm.CloseSession(sessionID)
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
			if n > 0 && onData != nil {
				encoded := base64.StdEncoding.EncodeToString(buf[:n])
				onData(sessionID, encoded)
			}
			if readErr != nil {
				return
			}
		}
	}()

	return &TerminalStartResult{
		SessionID: sessionID,
		Shell:     resolvedShell,
	}, nil
}

// Write writes data to the active PTY stdin
func (tm *TerminalManager) Write(sessionID string, data string) error {
	tm.mu.RLock()
	session, exists := tm.sessions[sessionID]
	tm.mu.RUnlock()

	if !exists || session == nil {
		return errors.New("session not found")
	}

	_, err := session.Conn.Write([]byte(data))
	return err
}

// Resize updates the terminal PTY rows and columns
func (tm *TerminalManager) Resize(ctx context.Context, cli *client.Client, sessionID string, rows uint, cols uint) error {
	tm.mu.RLock()
	session, exists := tm.sessions[sessionID]
	tm.mu.RUnlock()

	if !exists || session == nil {
		return errors.New("session not found")
	}

	if rows == 0 || cols == 0 {
		return nil
	}

	return cli.ContainerExecResize(ctx, session.ExecID, container.ResizeOptions{
		Height: rows,
		Width:  cols,
	})
}

// CloseSession closes a single terminal session
func (tm *TerminalManager) CloseSession(sessionID string) {
	tm.mu.Lock()
	session, exists := tm.sessions[sessionID]
	if exists {
		delete(tm.sessions, sessionID)
	}
	tm.mu.Unlock()

	if exists && session != nil {
		session.Close()
	}
}

// CloseAll terminates all active sessions
func (tm *TerminalManager) CloseAll() {
	tm.mu.Lock()
	all := make([]*TerminalSession, 0, len(tm.sessions))
	for _, s := range tm.sessions {
		all = append(all, s)
	}
	tm.sessions = make(map[string]*TerminalSession)
	tm.mu.Unlock()

	for _, s := range all {
		s.Close()
	}
}
