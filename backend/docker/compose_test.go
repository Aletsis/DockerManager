package docker

import (
	"context"
	"testing"
	"time"
)

func TestComposeStackOperationsValidation(t *testing.T) {
	svc, err := NewService()
	if err != nil {
		t.Skipf("Docker daemon not available: %v", err)
		return
	}
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Validation test: empty project name
	if err := svc.StartStack(ctx, ""); err == nil {
		t.Errorf("Expected error for empty project name on StartStack, got nil")
	}

	if err := svc.StopStack(ctx, ""); err == nil {
		t.Errorf("Expected error for empty project name on StopStack, got nil")
	}

	if err := svc.RestartStack(ctx, "   "); err == nil {
		t.Errorf("Expected error for whitespace project name on RestartStack, got nil")
	}

	// 2. Non-existent stack should return a clear not found error
	nonExistent := "non_existent_compose_project_12345"
	if err := svc.StartStack(ctx, nonExistent); err == nil {
		t.Errorf("Expected error when starting non-existent stack, got nil")
	}

	if err := svc.StopStack(ctx, nonExistent); err == nil {
		t.Errorf("Expected error when stopping non-existent stack, got nil")
	}

	if err := svc.RestartStack(ctx, nonExistent); err == nil {
		t.Errorf("Expected error when restarting non-existent stack, got nil")
	}
}
