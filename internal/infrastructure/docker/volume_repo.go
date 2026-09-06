package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	volumedomain "dockermanager/internal/domain/volume"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
)

// VolumeRepository implements volumedomain.Repository using Docker SDK
type VolumeRepository struct {
	client *Client
}

// NewVolumeRepository creates a new Docker VolumeRepository
func NewVolumeRepository(client *Client) *VolumeRepository {
	return &VolumeRepository{client: client}
}

// List returns all docker volumes enriched with attached container details and size
func (r *VolumeRepository) List(ctx context.Context) ([]volumedomain.Volume, error) {
	// 1. Fetch all containers to discover volume mounts
	containers, err := r.client.cli.ContainerList(ctx, container.ListOptions{All: true})
	volumeContainers := make(map[string][]volumedomain.ContainerRef)
	if err == nil {
		for _, c := range containers {
			cName := ""
			if len(c.Names) > 0 {
				cName = strings.TrimPrefix(c.Names[0], "/")
			}
			cID := c.ID
			if len(cID) > 12 {
				cID = cID[:12]
			}

			for _, m := range c.Mounts {
				if m.Type == "volume" || (m.Name != "" && m.Type != "bind") {
					volumeContainers[m.Name] = append(volumeContainers[m.Name], volumedomain.ContainerRef{
						ID:          cID,
						Name:        cName,
						State:       c.State,
						Destination: m.Destination,
						RW:          m.RW,
					})
				}
			}
		}
	}

	// 2. Fetch disk usage for size information (GET /system/df)
	volumeSizes := make(map[string]int64)
	du, err := r.client.cli.DiskUsage(ctx, types.DiskUsageOptions{})
	if err == nil {
		for _, v := range du.Volumes {
			if v != nil && v.UsageData != nil && v.UsageData.Size >= 0 {
				volumeSizes[v.Name] = v.UsageData.Size
			}
		}
	}

	// 3. Fetch all volumes
	volList, err := r.client.cli.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes: %w", err)
	}

	result := make([]volumedomain.Volume, 0, len(volList.Volumes))
	for _, v := range volList.Volumes {
		if v == nil {
			continue
		}

		refs := volumeContainers[v.Name]
		if refs == nil {
			refs = make([]volumedomain.ContainerRef, 0)
		}
		inUse := len(refs) > 0

		size := int64(-1)
		if sz, ok := volumeSizes[v.Name]; ok {
			size = sz
		} else if v.UsageData != nil && v.UsageData.Size >= 0 {
			size = v.UsageData.Size
		}

		result = append(result, volumedomain.Volume{
			Name:       v.Name,
			Driver:     v.Driver,
			Mountpoint: v.Mountpoint,
			CreatedAt:  v.CreatedAt,
			Labels:     v.Labels,
			Scope:      v.Scope,
			Size:       size,
			InUse:      inUse,
			Containers: refs,
		})
	}

	// Sort volumes: in-use first, then by name alphabetically
	sort.Slice(result, func(i, j int) bool {
		if result[i].InUse != result[j].InUse {
			return result[i].InUse && !result[j].InUse
		}
		return result[i].Name < result[j].Name
	})

	return result, nil
}

// GetDiskUsage retrieves aggregated volume disk metrics
func (r *VolumeRepository) GetDiskUsage(ctx context.Context) (*volumedomain.DiskUsageSummary, error) {
	vols, err := r.List(ctx)
	if err != nil {
		return nil, err
	}

	var totalSize, danglingSize, reclaimableSize int64
	danglingCount := 0

	for _, v := range vols {
		if v.Size > 0 {
			totalSize += v.Size
		}
		if !v.InUse {
			danglingCount++
			if v.Size > 0 {
				danglingSize += v.Size
				reclaimableSize += v.Size
			}
		}
	}

	return &volumedomain.DiskUsageSummary{
		TotalVolumes:    len(vols),
		TotalSize:       totalSize,
		DanglingCount:   danglingCount,
		DanglingSize:    danglingSize,
		ReclaimableSize: reclaimableSize,
	}, nil
}

// Inspect returns raw formatted JSON configuration of a volume
func (r *VolumeRepository) Inspect(ctx context.Context, name string) (string, error) {
	vol, err := r.client.cli.VolumeInspect(ctx, name)
	if err != nil {
		return "", fmt.Errorf("volume not found: %w", err)
	}

	bytes, err := json.MarshalIndent(vol, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format volume inspect json: %w", err)
	}

	return string(bytes), nil
}

// Remove deletes a volume
func (r *VolumeRepository) Remove(ctx context.Context, name string, force bool) error {
	err := r.client.cli.VolumeRemove(ctx, name, force)
	if err != nil {
		return fmt.Errorf("failed to remove volume %s: %w", name, err)
	}
	return nil
}

// Prune cleans all unused/dangling volumes
func (r *VolumeRepository) Prune(ctx context.Context) (*volumedomain.PruneResult, error) {
	report, err := r.client.cli.VolumesPrune(ctx, filters.NewArgs())
	if err != nil {
		return nil, fmt.Errorf("failed to prune volumes: %w", err)
	}

	return &volumedomain.PruneResult{
		VolumesDeleted: report.VolumesDeleted,
		SpaceReclaimed: report.SpaceReclaimed,
	}, nil
}
