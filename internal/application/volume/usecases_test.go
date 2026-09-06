package volume_test

import (
	"context"
	"testing"

	volumeapp "dockermanager/internal/application/volume"
	volumedomain "dockermanager/internal/domain/volume"
)

type mockVolumeRepo struct {
	volumes   []volumedomain.Volume
	diskUsage *volumedomain.DiskUsageSummary
	pruneRes  *volumedomain.PruneResult
	inspect   string
	err       error

	lastAction string
	lastName   string
	lastForce  bool
}

func (m *mockVolumeRepo) List(ctx context.Context) ([]volumedomain.Volume, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.volumes, nil
}

func (m *mockVolumeRepo) GetDiskUsage(ctx context.Context) (*volumedomain.DiskUsageSummary, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.diskUsage, nil
}

func (m *mockVolumeRepo) Inspect(ctx context.Context, name string) (string, error) {
	m.lastAction = "inspect"
	m.lastName = name
	return m.inspect, m.err
}

func (m *mockVolumeRepo) Remove(ctx context.Context, name string, force bool) error {
	m.lastAction = "remove"
	m.lastName = name
	m.lastForce = force
	return m.err
}

func (m *mockVolumeRepo) Prune(ctx context.Context) (*volumedomain.PruneResult, error) {
	m.lastAction = "prune"
	return m.pruneRes, m.err
}

func TestVolumeUseCases(t *testing.T) {
	mock := &mockVolumeRepo{
		volumes: []volumedomain.Volume{
			{Name: "data_vol", Driver: "local", Size: 1048576, InUse: true},
		},
		diskUsage: &volumedomain.DiskUsageSummary{
			TotalVolumes:    1,
			TotalSize:       1048576,
			DanglingCount:   0,
			DanglingSize:    0,
			ReclaimableSize: 0,
		},
		inspect: `{"Name": "data_vol"}`,
		pruneRes: &volumedomain.PruneResult{
			VolumesDeleted: []string{"orphan_vol"},
			SpaceReclaimed: 524288,
		},
	}

	ctx := context.Background()

	// 1. List
	listUC := volumeapp.NewListVolumesUseCase(mock)
	vols, err := listUC.Execute(ctx)
	if err != nil || len(vols) != 1 || vols[0].Name != "data_vol" {
		t.Fatalf("ListVolumesUseCase failed: %v", err)
	}

	// 2. DiskUsage
	duUC := volumeapp.NewGetVolumeDiskUsageUseCase(mock)
	du, err := duUC.Execute(ctx)
	if err != nil || du.TotalVolumes != 1 {
		t.Fatalf("GetVolumeDiskUsageUseCase failed: %v", err)
	}

	// 3. Inspect
	inspectUC := volumeapp.NewInspectVolumeUseCase(mock)
	jsonStr, err := inspectUC.Execute(ctx, "data_vol")
	if err != nil || jsonStr != `{"Name": "data_vol"}` || mock.lastName != "data_vol" {
		t.Fatalf("InspectVolumeUseCase failed: %v", err)
	}

	// 4. Remove
	removeUC := volumeapp.NewRemoveVolumeUseCase(mock)
	err = removeUC.Execute(ctx, "data_vol", true)
	if err != nil || mock.lastName != "data_vol" || !mock.lastForce {
		t.Fatalf("RemoveVolumeUseCase failed: %v", err)
	}

	// 5. Prune
	pruneUC := volumeapp.NewPruneVolumesUseCase(mock)
	res, err := pruneUC.Execute(ctx)
	if err != nil || len(res.VolumesDeleted) != 1 || res.SpaceReclaimed != 524288 {
		t.Fatalf("PruneVolumesUseCase failed: %v", err)
	}
}
