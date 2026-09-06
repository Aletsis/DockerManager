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
	}

	if !c.IsRunning() {
		t.Errorf("expected container to be running")
	}

	// CanStart should fail when already running
	if err := c.CanStart(); err != domain.ErrContainerAlreadyRunning {
		t.Errorf("expected ErrContainerAlreadyRunning, got %v", err)
	}

	// CanStop should succeed when running
	if err := c.CanStop(); err != nil {
		t.Errorf("expected CanStop to succeed, got %v", err)
	}

	// CanRemove should fail without force when running
	if err := c.CanRemove(false); err != domain.ErrCannotRemoveRunning {
		t.Errorf("expected ErrCannotRemoveRunning, got %v", err)
	}

	// CanRemove should succeed with force when running
	if err := c.CanRemove(true); err != nil {
		t.Errorf("expected CanRemove(true) to succeed, got %v", err)
	}

	// Switch state to exited
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
}

func TestValidateCreateSpec(t *testing.T) {
	err := container.ValidateCreateSpec(container.CreateSpec{Image: ""})
	if err != domain.ErrInvalidContainerSpec {
		t.Errorf("expected ErrInvalidContainerSpec, got %v", err)
	}

	err = container.ValidateCreateSpec(container.CreateSpec{Image: "nginx:alpine"})
	if err != nil {
		t.Errorf("expected valid spec to pass, got %v", err)
	}
}
