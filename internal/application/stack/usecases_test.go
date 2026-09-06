package stack_test

import (
	"context"
	"testing"

	stackapp "dockermanager/internal/application/stack"
)

type mockStackRepo struct {
	lastAction string
	lastStack  string
	err        error
}

func (m *mockStackRepo) StartStack(ctx context.Context, projectName string) error {
	m.lastAction = "start"
	m.lastStack = projectName
	return m.err
}

func (m *mockStackRepo) StopStack(ctx context.Context, projectName string) error {
	m.lastAction = "stop"
	m.lastStack = projectName
	return m.err
}

func (m *mockStackRepo) RestartStack(ctx context.Context, projectName string) error {
	m.lastAction = "restart"
	m.lastStack = projectName
	return m.err
}

func TestManageStackUseCase(t *testing.T) {
	mock := &mockStackRepo{}
	uc := stackapp.NewManageStackUseCase(mock)
	ctx := context.Background()

	_ = uc.StartStack(ctx, "app_stack")
	if mock.lastAction != "start" || mock.lastStack != "app_stack" {
		t.Fatalf("expected start app_stack, got %s %s", mock.lastAction, mock.lastStack)
	}

	_ = uc.StopStack(ctx, "app_stack")
	if mock.lastAction != "stop" || mock.lastStack != "app_stack" {
		t.Fatalf("expected stop app_stack, got %s %s", mock.lastAction, mock.lastStack)
	}

	_ = uc.RestartStack(ctx, "app_stack")
	if mock.lastAction != "restart" || mock.lastStack != "app_stack" {
		t.Fatalf("expected restart app_stack, got %s %s", mock.lastAction, mock.lastStack)
	}
}
