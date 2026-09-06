package integration_test

import (
	"context"
	"testing"
	"time"

	dockerinfra "dockermanager/internal/infrastructure/docker"
)

// getTestDockerClient returns a connected Docker client or skips the test if daemon is unreachable
func getTestDockerClient(t *testing.T) *dockerinfra.Client {
	t.Helper()

	client, err := dockerinfra.NewClient()
	if err != nil {
		t.Skipf("Skipping integration test: Docker client initialization failed: %v", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err = client.RawClient().Ping(ctx)
	if err != nil {
		t.Skipf("Skipping integration test: Docker daemon unreachable: %v", err)
		return nil
	}

	return client
}
