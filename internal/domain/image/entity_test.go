package image_test

import (
	"testing"

	"dockermanager/internal/domain"
	"dockermanager/internal/domain/image"
)

func TestImageInvariants(t *testing.T) {
	img := image.Image{
		ID:         "sha256:1234567890abcdef",
		ShortID:    "1234567890ab",
		Repository: "nginx",
		Tag:        "alpine",
		InUse:      true,
		IsDangling: false,
	}

	// When image is in use and not forced, CanRemove should fail
	if err := img.CanRemove(false); err != domain.ErrImageInUse {
		t.Errorf("expected ErrImageInUse when in use, got %v", err)
	}

	// When forced, CanRemove should succeed
	if err := img.CanRemove(true); err != nil {
		t.Errorf("expected CanRemove(true) to succeed, got %v", err)
	}

	// When not in use, CanRemove without force should succeed
	img.InUse = false
	if err := img.CanRemove(false); err != nil {
		t.Errorf("expected CanRemove(false) to succeed when not in use, got %v", err)
	}
}
