package integration_test

import (
	"context"
	"testing"
	"time"

	dockerinfra "dockermanager/internal/infrastructure/docker"
)

func TestImageIntegration(t *testing.T) {
	client := getTestDockerClient(t)
	repo := dockerinfra.NewImageRepository(client)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. List images
	images, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("Failed to list images: %v", err)
	}
	t.Logf("Found %d local images", len(images))

	// 2. Get disk usage
	du, err := repo.GetDiskUsage(ctx)
	if err != nil {
		t.Fatalf("Failed to query disk usage: %v", err)
	}

	if du.TotalImages < len(images) {
		t.Errorf("Disk usage TotalImages (%d) is less than ListImages count (%d)", du.TotalImages, len(images))
	}
	t.Logf("Disk usage: TotalImages=%d, TotalSize=%d bytes, ReclaimableSize=%d bytes",
		du.TotalImages, du.TotalSize, du.ReclaimableSize)
}
