package container_test

import (
	"testing"

	"dockermanager/internal/domain"
	"dockermanager/internal/domain/container"
)

func TestContainerInvariants(t *testing.T) {
	c := container.Container{
		ID:    "abc123456789",
		Name:  "test-container",
		State: container.StateRunning,
		Networks: []container.NetworkInfo{
			{NetworkName: "my-bridge", IPAddress: "172.18.0.2"},
		},
		ComposeProject: "my-stack",
	}

	if !c.IsRunning() {
		t.Errorf("expected container to be running")
	}
	if c.IsPaused() {
		t.Errorf("expected container not to be paused")
	}

	// CanStart should fail when already running
	if err := c.CanStart(); err != domain.ErrContainerAlreadyRunning {
		t.Errorf("expected ErrContainerAlreadyRunning, got %v", err)
	}

	// CanStop should succeed when running
	if err := c.CanStop(); err != nil {
		t.Errorf("expected CanStop to succeed, got %v", err)
	}

	// CanPause should succeed when running
	if err := c.CanPause(); err != nil {
		t.Errorf("expected CanPause to succeed, got %v", err)
	}

	// CanUnpause should fail when running (not paused)
	if err := c.CanUnpause(); err != domain.ErrContainerNotPaused {
		t.Errorf("expected ErrContainerNotPaused, got %v", err)
	}

	// CanRemove should fail without force when running
	if err := c.CanRemove(false); err != domain.ErrCannotRemoveRunning {
		t.Errorf("expected ErrCannotRemoveRunning, got %v", err)
	}

	// CanRemove should succeed with force when running
	if err := c.CanRemove(true); err != nil {
		t.Errorf("expected CanRemove(true) to succeed, got %v", err)
	}

	// Network & Compose checks
	if !c.BelongsToCompose() {
		t.Errorf("expected container to belong to compose")
	}
	if !c.BelongsToNetwork("my-bridge") {
		t.Errorf("expected container to belong to my-bridge network")
	}
	if c.BelongsToNetwork("other-net") {
		t.Errorf("expected container not to belong to other-net")
	}

	// Switch state to Paused
	c.State = container.StatePaused
	if !c.IsPaused() {
		t.Errorf("expected container to be paused")
	}
	if err := c.CanStart(); err != domain.ErrContainerPaused {
		t.Errorf("expected ErrContainerPaused when trying to start paused, got %v", err)
	}
	if err := c.CanPause(); err != domain.ErrContainerNotRunning {
		t.Errorf("expected ErrContainerNotRunning when trying to pause paused, got %v", err)
	}
	if err := c.CanUnpause(); err != nil {
		t.Errorf("expected CanUnpause to succeed for paused container, got %v", err)
	}

	// Switch state to Exited
	c.State = container.StateExited
	if c.IsRunning() {
		t.Errorf("expected container not to be running")
	}
	if err := c.CanStart(); err != nil {
		t.Errorf("expected CanStart to succeed for exited container, got %v", err)
	}
	if err := c.CanStop(); err != domain.ErrContainerNotRunning {
		t.Errorf("expected ErrContainerNotRunning, got %v", err)
	}
	if err := c.CanRemove(false); err != nil {
		t.Errorf("expected CanRemove(false) to succeed for exited container, got %v", err)
	}
}

func TestValidateCreateSpec(t *testing.T) {
	err := container.ValidateCreateSpec(container.CreateSpec{Image: ""})
	if err != domain.ErrInvalidContainerSpec {
		t.Errorf("expected ErrInvalidContainerSpec, got %v", err)
	}

	err = container.ValidateCreateSpec(container.CreateSpec{Image: "   "})
	if err != domain.ErrInvalidContainerSpec {
		t.Errorf("expected ErrInvalidContainerSpec for whitespace, got %v", err)
	}

	err = container.ValidateCreateSpec(container.CreateSpec{Image: "nginx:alpine"})
	if err != nil {
		t.Errorf("expected valid spec to pass, got %v", err)
	}
}
