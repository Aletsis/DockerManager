package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/jsonmessage"
)

// ListImages returns list of all local images with structured metadata and usage status
func (s *Service) ListImages(ctx context.Context) ([]ImageInfo, error) {
	images, err := s.cli.ImageList(ctx, image.ListOptions{All: false})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	// Cross reference running/stopped containers to find images in use
	containers, _ := s.cli.ContainerList(ctx, container.ListOptions{All: true})
	usedImages := make(map[string]bool)
	for _, c := range containers {
		usedImages[c.ImageID] = true
		usedImages[c.Image] = true
	}

	result := make([]ImageInfo, 0, len(images))
	for _, img := range images {
		shortID := img.ID
		if strings.HasPrefix(shortID, "sha256:") {
			shortID = strings.TrimPrefix(shortID, "sha256:")
		}
		if len(shortID) > 12 {
			shortID = shortID[:12]
		}

		var firstTag string
		if len(img.RepoTags) > 0 {
			firstTag = img.RepoTags[0]
		}
		repo, tag, isDangling := ParseRepoTag(firstTag)

		inUse := img.Containers > 0 || usedImages[img.ID]
		if !inUse {
			for _, rt := range img.RepoTags {
				if usedImages[rt] {
					inUse = true
					break
				}
			}
		}

		result = append(result, ImageInfo{
			ID:         img.ID,
			ShortID:    shortID,
			Repository: repo,
			Tag:        tag,
			RepoTags:   img.RepoTags,
			Created:    img.Created,
			Size:       img.Size,
			SharedSize: img.SharedSize,
			Containers: img.Containers,
			InUse:      inUse,
			IsDangling: isDangling,
		})
	}
	return result, nil
}

// GetDiskUsage returns aggregated metrics of image disk storage and recoverable space
func (s *Service) GetDiskUsage(ctx context.Context) (*DiskUsageSummary, error) {
	du, err := s.cli.DiskUsage(ctx, types.DiskUsageOptions{})
	if err != nil {
		// Fallback: calculate directly from ListImages if DiskUsage endpoint is restricted
		images, listErr := s.ListImages(ctx)
		if listErr != nil {
			return nil, fmt.Errorf("failed to retrieve disk usage: %w", err)
		}
		var totalSize, danglingSize, reclaimableSize int64
		danglingCount := 0
		for _, img := range images {
			totalSize += img.Size
			if img.IsDangling {
				danglingCount++
				danglingSize += img.Size
			}
			if !img.InUse {
				reclaimableSize += img.Size
			}
		}
		return &DiskUsageSummary{
			TotalImages:     len(images),
			TotalSize:       totalSize,
			DanglingCount:   danglingCount,
			DanglingSize:    danglingSize,
			ReclaimableSize: reclaimableSize,
		}, nil
	}

	var totalSize, danglingSize, reclaimableSize int64
	danglingCount := 0

	containers, _ := s.cli.ContainerList(ctx, container.ListOptions{All: true})
	usedImages := make(map[string]bool)
	for _, c := range containers {
		usedImages[c.ImageID] = true
		usedImages[c.Image] = true
	}

	for _, img := range du.Images {
		totalSize += img.Size
		var firstTag string
		if len(img.RepoTags) > 0 {
			firstTag = img.RepoTags[0]
		}
		_, _, isDangling := ParseRepoTag(firstTag)
		if isDangling {
			danglingCount++
			danglingSize += img.Size
		}
		inUse := img.Containers > 0 || usedImages[img.ID]
		if !inUse {
			for _, rt := range img.RepoTags {
				if usedImages[rt] {
					inUse = true
					break
				}
			}
		}
		if !inUse {
			reclaimableSize += img.Size
		}
	}

	return &DiskUsageSummary{
		TotalImages:     len(du.Images),
		TotalSize:       totalSize,
		DanglingCount:   danglingCount,
		DanglingSize:    danglingSize,
		ReclaimableSize: reclaimableSize,
	}, nil
}

// PullImage downloads an image while streaming progress events
func (s *Service) PullImage(ctx context.Context, ref string, onProgress func(event PullProgressEvent)) error {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return fmt.Errorf("el nombre de la imagen no puede estar vacío")
	}

	reader, err := s.cli.ImagePull(ctx, ref, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("error al iniciar descarga de imagen %s: %w", ref, err)
	}
	defer reader.Close()

	decoder := json.NewDecoder(reader)
	for {
		var msg jsonmessage.JSONMessage
		if err := decoder.Decode(&msg); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error al leer respuesta de descarga: %w", err)
		}

		if msg.Error != nil {
			if onProgress != nil {
				onProgress(PullProgressEvent{
					ID:     msg.ID,
					Status: "Error",
					Error:  msg.Error.Message,
				})
			}
			return fmt.Errorf("error descargando imagen: %s", msg.Error.Message)
		}

		if onProgress != nil {
			var cur, tot int64
			if msg.Progress != nil {
				cur = msg.Progress.Current
				tot = msg.Progress.Total
			}
			onProgress(PullProgressEvent{
				ID:       msg.ID,
				Status:   msg.Status,
				Progress: msg.ProgressMessage,
				Current:  cur,
				Total:    tot,
			})
		}
	}

	return nil
}

// RemoveImage deletes an image by ID or name
func (s *Service) RemoveImage(ctx context.Context, id string, force bool) error {
	_, err := s.cli.ImageRemove(ctx, id, image.RemoveOptions{
		Force:         force,
		PruneChildren: true,
	})
	if err != nil {
		return fmt.Errorf("error al eliminar imagen %s: %w", id, err)
	}
	return nil
}

// PruneImages deletes unused or dangling images and returns reclaimed space
func (s *Service) PruneImages(ctx context.Context, danglingOnly bool) (*PruneResult, error) {
	pruneFilters := filters.NewArgs()
	if danglingOnly {
		pruneFilters.Add("dangling", "true")
	} else {
		pruneFilters.Add("dangling", "false")
	}

	report, err := s.cli.ImagesPrune(ctx, pruneFilters)
	if err != nil {
		return nil, fmt.Errorf("error al ejecutar limpieza de imágenes: %w", err)
	}

	deletedIDs := make([]string, 0, len(report.ImagesDeleted))
	for _, item := range report.ImagesDeleted {
		if item.Deleted != "" {
			deletedIDs = append(deletedIDs, item.Deleted)
		} else if item.Untagged != "" {
			deletedIDs = append(deletedIDs, item.Untagged)
		}
	}

	return &PruneResult{
		ImagesDeleted:  deletedIDs,
		SpaceReclaimed: report.SpaceReclaimed,
	}, nil
}
