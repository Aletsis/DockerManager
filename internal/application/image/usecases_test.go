package image_test

import (
	"context"
	"testing"

	imageapp "dockermanager/internal/application/image"
	imagedomain "dockermanager/internal/domain/image"
)

type mockImageRepo struct {
	images    []imagedomain.Image
	diskUsage *imagedomain.DiskUsageSummary
	pruneRes  *imagedomain.PruneResult
	err       error

	lastAction string
	lastRef    string
	lastID     string
	lastForce  bool
}

func (m *mockImageRepo) List(ctx context.Context) ([]imagedomain.Image, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.images, nil
}

func (m *mockImageRepo) GetDiskUsage(ctx context.Context) (*imagedomain.DiskUsageSummary, error) {
	return m.diskUsage, m.err
}

func (m *mockImageRepo) Pull(ctx context.Context, imageName string, onProgress func(imagedomain.PullProgressEvent)) error {
	m.lastAction = "pull"
	m.lastRef = imageName
	if onProgress != nil {
		onProgress(imagedomain.PullProgressEvent{ID: "layer1", Status: "Downloading", Current: 50, Total: 100})
	}
	return m.err
}

func (m *mockImageRepo) Remove(ctx context.Context, id string, force bool) error {
	m.lastAction = "remove"
	m.lastID = id
	m.lastForce = force
	return m.err
}

func (m *mockImageRepo) Prune(ctx context.Context, danglingOnly bool) (*imagedomain.PruneResult, error) {
	m.lastAction = "prune"
	return m.pruneRes, m.err
}

func TestImageUseCases(t *testing.T) {
	mock := &mockImageRepo{
		images: []imagedomain.Image{
			{ID: "img1", Repository: "nginx", Tag: "alpine"},
		},
		diskUsage: &imagedomain.DiskUsageSummary{TotalImages: 1, TotalSize: 50000000},
		pruneRes:  &imagedomain.PruneResult{ImagesDeleted: []string{"img_old"}, SpaceReclaimed: 1000000},
	}
	ctx := context.Background()

	// List
	listUC := imageapp.NewListImagesUseCase(mock)
	images, err := listUC.Execute(ctx)
	if err != nil || len(images) != 1 {
		t.Fatalf("expected 1 image, got %d, err: %v", len(images), err)
	}

	// Disk Usage
	diskUC := imageapp.NewGetDiskUsageUseCase(mock)
	du, err := diskUC.Execute(ctx)
	if err != nil || du.TotalImages != 1 {
		t.Fatalf("unexpected disk usage: %v, err: %v", du, err)
	}

	// Pull
	pullUC := imageapp.NewPullImageUseCase(mock)
	var progressReceived bool
	err = pullUC.Execute(ctx, "redis:alpine", func(e imagedomain.PullProgressEvent) {
		progressReceived = true
	})
	if err != nil || mock.lastRef != "redis:alpine" || !progressReceived {
		t.Fatalf("pull failed or didn't receive progress callback")
	}

	// Remove
	removeUC := imageapp.NewRemoveImageUseCase(mock)
	err = removeUC.Execute(ctx, "img1", true)
	if err != nil || mock.lastID != "img1" || !mock.lastForce {
		t.Fatalf("remove failed: id=%s force=%v", mock.lastID, mock.lastForce)
	}

	// Prune
	pruneUC := imageapp.NewPruneImagesUseCase(mock)
	res, err := pruneUC.Execute(ctx, true)
	if err != nil || len(res.ImagesDeleted) != 1 || res.SpaceReclaimed != 1000000 {
		t.Fatalf("prune failed: %v", res)
	}
}
