package docker

import (
	"context"
	"testing"
	"time"
)

func TestNetworkOperationsValidation(t *testing.T) {
	svc, err := NewService()
	if err != nil {
		t.Skipf("Docker daemon not available: %v", err)
		return
	}
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Validation test: empty network name
	if err := svc.StartNetwork(ctx, ""); err == nil {
		t.Errorf("Expected error for empty network name on StartNetwork, got nil")
	}

	if err := svc.StopNetwork(ctx, ""); err == nil {
		t.Errorf("Expected error for empty network name on StopNetwork, got nil")
	}

	if err := svc.RestartNetwork(ctx, "   "); err == nil {
		t.Errorf("Expected error for whitespace network name on RestartNetwork, got nil")
	}

	// 2. Non-existent network should return a clear not found error
	nonExistent := "non_existent_docker_network_xyz_98765"
	if err := svc.StartNetwork(ctx, nonExistent); err == nil {
		t.Errorf("Expected error when starting non-existent network, got nil")
	}

	if err := svc.StopNetwork(ctx, nonExistent); err == nil {
		t.Errorf("Expected error when stopping non-existent network, got nil")
	}

	if err := svc.RestartNetwork(ctx, nonExistent); err == nil {
		t.Errorf("Expected error when restarting non-existent network, got nil")
	}
}
