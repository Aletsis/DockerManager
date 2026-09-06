package integration_test

import (
	"context"
	"testing"
	"time"

	dockerinfra "dockermanager/internal/infrastructure/docker"
)

func TestSystemIntegration(t *testing.T) {
	client := getTestDockerClient(t)
	repo := dockerinfra.NewSystemRepository(client)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	overview, err := repo.GetOverview(ctx)
	if err != nil {
		t.Fatalf("Failed to retrieve system overview: %v", err)
	}

	if overview.NCPU <= 0 {
		t.Errorf("expected NCPU > 0, got %d", overview.NCPU)
	}
	if overview.OperatingSystem == "" {
		t.Errorf("expected non-empty OperatingSystem")
	}
	if overview.ServerVersion == "" {
		t.Errorf("expected non-empty ServerVersion")
	}

	t.Logf("Docker host: OS=%s, Version=%s, Containers=%d, Images=%d",
		overview.OperatingSystem, overview.ServerVersion, overview.Containers, overview.Images)
}
