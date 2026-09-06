package volume

import "dockermanager/internal/domain"

// ContainerRef represents a container referencing/mounting a volume
type ContainerRef struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	State       string `json:"state"`
	Destination string `json:"destination"`
	RW          bool   `json:"rw"`
}

// Volume represents a Docker storage volume
type Volume struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Mountpoint string            `json:"mountpoint"`
	CreatedAt  string            `json:"createdAt"`
	Labels     map[string]string `json:"labels"`
	Scope      string            `json:"scope"`
	Size       int64             `json:"size"` // estimated size in bytes, -1 if unavailable
	InUse      bool              `json:"inUse"`
	Containers []ContainerRef    `json:"containers"`
}

// CanRemove checks whether a volume can be safely removed
func (v *Volume) CanRemove(force bool) error {
	if v.InUse && !force {
		return domain.ErrVolumeInUse
	}
	return nil
}

// DiskUsageSummary provides aggregated disk metrics for volumes
type DiskUsageSummary struct {
	TotalVolumes    int   `json:"totalVolumes"`
	TotalSize       int64 `json:"totalSize"`
	DanglingCount   int   `json:"danglingCount"`
	DanglingSize    int64 `json:"danglingSize"`
	ReclaimableSize int64 `json:"reclaimableSize"`
}

// PruneResult contains the result of volume pruning
type PruneResult struct {
	VolumesDeleted []string `json:"volumesDeleted"`
	SpaceReclaimed uint64   `json:"spaceReclaimed"`
}
