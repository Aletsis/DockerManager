package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	imagedomain "dockermanager/internal/domain/image"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/jsonmessage"
)

// ImageRepository implements imagedomain.Repository using Docker SDK
type ImageRepository struct {
	client *Client
}

// NewImageRepository creates a new Docker ImageRepository
func NewImageRepository(client *Client) *ImageRepository {
	return &ImageRepository{client: client}
}

// List returns all local images with domain mapping
func (r *ImageRepository) List(ctx context.Context) ([]imagedomain.Image, error) {
	images, err := r.client.cli.ImageList(ctx, image.ListOptions{All: false})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	containers, _ := r.client.cli.ContainerList(ctx, container.ListOptions{All: true})
	usedImages := make(map[string]bool)
	for _, c := range containers {
		usedImages[c.ImageID] = true
		usedImages[c.Image] = true
	}

	result := make([]imagedomain.Image, 0, len(images))
	for _, img := range images {
		inUse := img.Containers > 0 || usedImages[img.ID]
		if !inUse {
			for _, rt := range img.RepoTags {
				if usedImages[rt] {
					inUse = true
					break
				}
			}
		}
		result = append(result, toDomainImage(img, inUse))
	}
	return result, nil
}

// GetDiskUsage calculates disk usage metrics
func (r *ImageRepository) GetDiskUsage(ctx context.Context) (*imagedomain.DiskUsageSummary, error) {
	du, err := r.client.cli.DiskUsage(ctx, types.DiskUsageOptions{})
	if err != nil {
		images, listErr := r.List(ctx)
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
		return &imagedomain.DiskUsageSummary{
			TotalImages:     len(images),
			TotalSize:       totalSize,
			DanglingCount:   danglingCount,
			DanglingSize:    danglingSize,
			ReclaimableSize: reclaimableSize,
		}, nil
	}

	var totalSize, danglingSize, reclaimableSize int64
	danglingCount := 0

	containers, _ := r.client.cli.ContainerList(ctx, container.ListOptions{All: true})
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

	return &imagedomain.DiskUsageSummary{
		TotalImages:     len(du.Images),
		TotalSize:       totalSize,
		DanglingCount:   danglingCount,
		DanglingSize:    danglingSize,
		ReclaimableSize: reclaimableSize,
	}, nil
}

// Pull downloads an image and streams progress events
func (r *ImageRepository) Pull(ctx context.Context, ref string, onProgress func(event imagedomain.PullProgressEvent)) error {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return fmt.Errorf("el nombre de la imagen no puede estar vacío")
	}

	reader, err := r.client.cli.ImagePull(ctx, ref, image.PullOptions{})
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
				onProgress(imagedomain.PullProgressEvent{
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
			onProgress(imagedomain.PullProgressEvent{
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

// Remove deletes an image
func (r *ImageRepository) Remove(ctx context.Context, id string, force bool) error {
	_, err := r.client.cli.ImageRemove(ctx, id, image.RemoveOptions{
		Force:         force,
		PruneChildren: true,
	})
	if err != nil {
		return fmt.Errorf("error al eliminar imagen %s: %w", id, err)
	}
	return nil
}

// Prune cleans unused images
func (r *ImageRepository) Prune(ctx context.Context, danglingOnly bool) (*imagedomain.PruneResult, error) {
	pruneFilters := filters.NewArgs()
	if danglingOnly {
		pruneFilters.Add("dangling", "true")
	} else {
		pruneFilters.Add("dangling", "false")
	}

	report, err := r.client.cli.ImagesPrune(ctx, pruneFilters)
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

	return &imagedomain.PruneResult{
		ImagesDeleted:  deletedIDs,
		SpaceReclaimed: report.SpaceReclaimed,
	}, nil
}
