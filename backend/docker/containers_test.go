package docker

import (
	"context"
	"testing"
	"time"
)

func TestCreateContainerService(t *testing.T) {
	svc, err := NewService()
	if err != nil {
		t.Skipf("Docker daemon not available: %v", err)
		return
	}
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 1. Validation test: empty image
	_, err = svc.CreateContainer(ctx, CreateContainerRequest{
		Image: "",
	})
	if err == nil {
		t.Errorf("Expected error for empty image name, got nil")
	}

	// 2. Integration test: create and clean up a test container
	testContainerName := "dockermanager-test-container"
	// Ensure no leftover from previous runs
	_ = svc.RemoveContainer(ctx, testContainerName, true)

	req := CreateContainerRequest{
		Image:         "alpine:latest",
		Name:          testContainerName,
		Ports:         []string{"18080:80/tcp"},
		Env:           []string{"TEST_ENV=antigravity"},
		RestartPolicy: "no",
		AutoStart:     false,
	}

	res, err := svc.CreateContainer(ctx, req)
	if err != nil {
		t.Fatalf("CreateContainer failed: %v", err)
	}

	if res.ID == "" {
		t.Errorf("Expected valid container ID, got empty string")
	}
	t.Logf("Created container successfully with ID: %s", res.ID)

	// Clean up test container
	err = svc.RemoveContainer(ctx, res.ID, true)
	if err != nil {
		t.Errorf("Failed to clean up test container %s: %v", res.ID, err)
	} else {
		t.Logf("Cleaned up container %s", res.ID)
	}
}
