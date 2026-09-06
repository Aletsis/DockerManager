package volume_test

import (
	"testing"

	"dockermanager/internal/domain"
	volumedomain "dockermanager/internal/domain/volume"
)

func TestVolumeCanRemove(t *testing.T) {
	inUseVol := volumedomain.Volume{
		Name:  "db_data",
		InUse: true,
		Containers: []volumedomain.ContainerRef{
			{ID: "c1", Name: "postgres_db"},
		},
	}

	// Should not allow removing in-use volume without force
	if err := inUseVol.CanRemove(false); err != domain.ErrVolumeInUse {
		t.Fatalf("expected ErrVolumeInUse, got %v", err)
	}

	// Should allow removing in-use volume with force
	if err := inUseVol.CanRemove(true); err != nil {
		t.Fatalf("expected nil when force is true, got %v", err)
	}

	unusedVol := volumedomain.Volume{
		Name:       "orphan_vol",
		InUse:      false,
		Containers: []volumedomain.ContainerRef{},
	}

	// Should allow removing unused volume without force
	if err := unusedVol.CanRemove(false); err != nil {
		t.Fatalf("expected nil for unused volume, got %v", err)
	}
}
