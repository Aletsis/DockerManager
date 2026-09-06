package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	containerdomain "dockermanager/internal/domain/container"
	dockerinfra "dockermanager/internal/infrastructure/docker"
)

func TestContainerIntegration(t *testing.T) {
	client := getTestDockerClient(t)
	repo := dockerinfra.NewContainerRepository(client)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. List existing containers
	existing, err := repo.List(ctx, true)
	if err != nil {
		t.Fatalf("Failed to list containers: %v", err)
	}
	t.Logf("Found %d existing containers", len(existing))

	// 2. Create a temporary integration test container
	testContainerName := fmt.Sprintf("dockermanager_test_%d", time.Now().UnixNano())
	createSpec := containerdomain.CreateSpec{
		Image:     "alpine:latest",
		Name:      testContainerName,
		AutoStart: false,
	}

	createRes, err := repo.Create(ctx, createSpec)
	if err != nil {
		t.Skipf("Skipping container creation test (could not pull/create alpine): %v", err)
		return
	}

	defer func() {
		_ = repo.Remove(context.Background(), createRes.ID, true)
		t.Logf("Cleaned up integration container %s", createRes.ID)
	}()

	// 3. Start the container
	err = repo.Start(ctx, createRes.ID)
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	// 4. Verify it is running
	found, err := repo.FindByID(ctx, createRes.ID)
	if err != nil {
		t.Fatalf("Failed to find created container: %v", err)
	}
	if !found.IsRunning() {
		t.Errorf("Expected container %s to be running, got state: %s", createRes.ID, found.State)
	}

	// 5. Stop the container
	err = repo.Stop(ctx, createRes.ID)
	if err != nil {
		t.Fatalf("Failed to stop container: %v", err)
	}
}
