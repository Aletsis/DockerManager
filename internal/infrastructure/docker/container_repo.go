package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	containerdomain "dockermanager/internal/domain/container"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/stdcopy"
)

// ContainerRepository implements containerdomain.Repository using Docker SDK
type ContainerRepository struct {
	client *Client
}

// NewContainerRepository creates a new Docker ContainerRepository
func NewContainerRepository(client *Client) *ContainerRepository {
	return &ContainerRepository{client: client}
}

// List returns a list of domain containers
func (r *ContainerRepository) List(ctx context.Context, all bool) ([]containerdomain.Container, error) {
	containers, err := r.client.cli.ContainerList(ctx, container.ListOptions{All: all})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	result := make([]containerdomain.Container, 0, len(containers))
	for _, c := range containers {
		result = append(result, toDomainContainer(c))
	}
	return result, nil
}

// FindByID retrieves a container by its ID
func (r *ContainerRepository) FindByID(ctx context.Context, id string) (*containerdomain.Container, error) {
	jsonContainer, err := r.client.cli.ContainerInspect(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("container not found: %w", err)
	}

	c := containerdomain.Container{
		ID:      jsonContainer.ID,
		Name:    strings.TrimPrefix(jsonContainer.Name, "/"),
		Image:   jsonContainer.Config.Image,
		ImageID: jsonContainer.Image,
		Created: 0,
		State:   containerdomain.State(jsonContainer.State.Status),
		Status:  jsonContainer.State.Status,
	}
	if len(c.ID) > 12 {
		c.ShortID = c.ID[:12]
	} else {
		c.ShortID = c.ID
	}
	return &c, nil
}

// Start starts a container
func (r *ContainerRepository) Start(ctx context.Context, id string) error {
	return r.client.cli.ContainerStart(ctx, id, container.StartOptions{})
}

// Stop stops a container
func (r *ContainerRepository) Stop(ctx context.Context, id string) error {
	timeout := 15
	return r.client.cli.ContainerStop(ctx, id, container.StopOptions{Timeout: &timeout})
}

// Restart restarts a container
func (r *ContainerRepository) Restart(ctx context.Context, id string) error {
	timeout := 15
	return r.client.cli.ContainerRestart(ctx, id, container.StopOptions{Timeout: &timeout})
}

// Pause pauses a container
func (r *ContainerRepository) Pause(ctx context.Context, id string) error {
	return r.client.cli.ContainerPause(ctx, id)
}

// Unpause unpauses a container
func (r *ContainerRepository) Unpause(ctx context.Context, id string) error {
	return r.client.cli.ContainerUnpause(ctx, id)
}

// Remove deletes a container
func (r *ContainerRepository) Remove(ctx context.Context, id string, force bool) error {
	return r.client.cli.ContainerRemove(ctx, id, container.RemoveOptions{
		Force:         force,
		RemoveVolumes: true,
	})
}

// GetLogs returns formatted container logs
func (r *ContainerRepository) GetLogs(ctx context.Context, id string, tail int) (string, error) {
	tailStr := strconv.Itoa(tail)
	if tail <= 0 {
		tailStr = "200"
	}

	reader, err := r.client.cli.ContainerLogs(ctx, id, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tailStr,
		Timestamps: false,
	})
	if err != nil {
		return "", fmt.Errorf("failed to fetch container logs: %w", err)
	}
	defer reader.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, reader)
	if err != nil {
		raw, readErr := io.ReadAll(reader)
		if readErr == nil && len(raw) > 0 {
			return string(raw), nil
		}
		return stdoutBuf.String() + stderrBuf.String(), nil
	}

	var combined strings.Builder
	if stdoutBuf.Len() > 0 {
		combined.WriteString(stdoutBuf.String())
	}
	if stderrBuf.Len() > 0 {
		combined.WriteString(stderrBuf.String())
	}

	return combined.String(), nil
}

// GetStats returns current container resource metrics
func (r *ContainerRepository) GetStats(ctx context.Context, id string) (*containerdomain.Stats, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	statsResponse, err := r.client.cli.ContainerStats(ctxTimeout, id, false)
	if err != nil {
		return nil, fmt.Errorf("failed to query stats: %w", err)
	}
	defer statsResponse.Body.Close()

	var stats container.StatsResponse
	decoder := json.NewDecoder(statsResponse.Body)
	if err := decoder.Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode stats: %w", err)
	}

	return CalculateContainerStats(stats.ID, stats.Name, &stats), nil
}

// Inspect returns formatted inspect JSON of the container
func (r *ContainerRepository) Inspect(ctx context.Context, id string) (string, error) {
	_, raw, err := r.client.cli.ContainerInspectWithRaw(ctx, id, false)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container: %w", err)
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, raw, "", "  "); err == nil {
		return prettyJSON.String(), nil
	}
	return string(raw), nil
}

// Create builds and creates a container based on the domain spec
func (r *ContainerRepository) Create(ctx context.Context, spec containerdomain.CreateSpec) (*containerdomain.CreateResult, error) {
	containerConfig, hostConfig, err := BuildContainerConfig(spec)
	if err != nil {
		return nil, err
	}

	// Verify if image exists locally; if not, pull it
	_, _, err = r.client.cli.ImageInspectWithRaw(ctx, spec.Image)
	if err != nil {
		reader, pullErr := r.client.cli.ImagePull(ctx, spec.Image, image.PullOptions{})
		if pullErr != nil {
			return nil, fmt.Errorf("la imagen %s no está disponible y falló la descarga: %w", spec.Image, pullErr)
		}
		_, _ = io.Copy(io.Discard, reader)
		_ = reader.Close()
	}

	createResp, err := r.client.cli.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		nil,
		nil,
		spec.Name,
	)
	if err != nil {
		return nil, fmt.Errorf("error al crear el contenedor: %w", err)
	}

	if spec.AutoStart {
		if startErr := r.client.cli.ContainerStart(ctx, createResp.ID, container.StartOptions{}); startErr != nil {
			return &containerdomain.CreateResult{
				ID:       createResp.ID,
				Warnings: createResp.Warnings,
			}, fmt.Errorf("contenedor creado con ID %s pero falló al iniciar: %w", createResp.ID, startErr)
		}
	}

	return &containerdomain.CreateResult{
		ID:       createResp.ID,
		Warnings: createResp.Warnings,
	}, nil
}
