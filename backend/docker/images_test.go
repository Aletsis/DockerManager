package docker

import (
	"context"
	"testing"
	"time"
)

func TestImagesService(t *testing.T) {
	svc, err := NewService()
	if err != nil {
		t.Skipf("Docker daemon not available: %v", err)
		return
	}
	defer svc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Test ListImages
	images, err := svc.ListImages(ctx)
	if err != nil {
		t.Fatalf("ListImages failed: %v", err)
	}
	t.Logf("Found %d images", len(images))

	hasInUse := false
	for _, img := range images {
		if img.InUse {
			hasInUse = true
			t.Logf("In-use image found: %s:%s (size: %d bytes)", img.Repository, img.Tag, img.Size)
			break
		}
	}
	t.Logf("Has in-use images: %v", hasInUse)

	// 2. Test GetDiskUsage
	du, err := svc.GetDiskUsage(ctx)
	if err != nil {
		t.Fatalf("GetDiskUsage failed: %v", err)
	}
	t.Logf("Disk usage: TotalSize=%d bytes, DanglingCount=%d, DanglingSize=%d bytes, ReclaimableSize=%d bytes",
		du.TotalSize, du.DanglingCount, du.DanglingSize, du.ReclaimableSize)

	if du.TotalImages != len(images) {
		t.Errorf("Expected TotalImages to be %d, got %d", len(images), du.TotalImages)
	}
}
