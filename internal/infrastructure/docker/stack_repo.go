package docker

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

// StackRepository implements stackdomain.Repository using Docker SDK
type StackRepository struct {
	client *Client
}

// NewStackRepository creates a new Docker StackRepository
func NewStackRepository(client *Client) *StackRepository {
	return &StackRepository{client: client}
}

// StartStack starts all stopped containers in a compose project
func (r *StackRepository) StartStack(ctx context.Context, projectName string) error {
	projectName = strings.TrimSpace(projectName)
	if projectName == "" {
		return fmt.Errorf("el nombre del proyecto compose no puede estar vacío")
	}

	filterArgs := filters.NewArgs()
	filterArgs.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))

	containers, err := r.client.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filterArgs})
	if err != nil {
		return fmt.Errorf("error al listar contenedores del stack %s: %w", projectName, err)
	}
	if len(containers) == 0 {
		return fmt.Errorf("no se encontraron contenedores para el stack %q", projectName)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []string

	for _, c := range containers {
		if c.State == "running" {
			continue
		}
		wg.Add(1)
		go func(cid string) {
			defer wg.Done()
			if err := r.client.cli.ContainerStart(ctx, cid, container.StartOptions{}); err != nil {
				mu.Lock()
				short := cid
				if len(short) > 12 {
					short = short[:12]
				}
				errs = append(errs, fmt.Sprintf("%s: %v", short, err))
				mu.Unlock()
			}
		}(c.ID)
	}
	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("errores al iniciar el stack: %s", strings.Join(errs, "; "))
	}
	return nil
}

// StopStack stops all running containers in a compose project
func (r *StackRepository) StopStack(ctx context.Context, projectName string) error {
	projectName = strings.TrimSpace(projectName)
	if projectName == "" {
		return fmt.Errorf("el nombre del proyecto compose no puede estar vacío")
	}

	filterArgs := filters.NewArgs()
	filterArgs.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))

	containers, err := r.client.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filterArgs})
	if err != nil {
		return fmt.Errorf("error al listar contenedores del stack %s: %w", projectName, err)
	}
	if len(containers) == 0 {
		return fmt.Errorf("no se encontraron contenedores para el stack %q", projectName)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []string
	timeout := 15

	for _, c := range containers {
		if c.State != "running" {
			continue
		}
		wg.Add(1)
		go func(cid string) {
			defer wg.Done()
			if err := r.client.cli.ContainerStop(ctx, cid, container.StopOptions{Timeout: &timeout}); err != nil {
				mu.Lock()
				short := cid
				if len(short) > 12 {
					short = short[:12]
				}
				errs = append(errs, fmt.Sprintf("%s: %v", short, err))
				mu.Unlock()
			}
		}(c.ID)
	}
	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("errores al detener el stack: %s", strings.Join(errs, "; "))
	}
	return nil
}

// RestartStack restarts all containers in a compose project
func (r *StackRepository) RestartStack(ctx context.Context, projectName string) error {
	projectName = strings.TrimSpace(projectName)
	if projectName == "" {
		return fmt.Errorf("el nombre del proyecto compose no puede estar vacío")
	}

	filterArgs := filters.NewArgs()
	filterArgs.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))

	containers, err := r.client.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filterArgs})
	if err != nil {
		return fmt.Errorf("error al listar contenedores del stack %s: %w", projectName, err)
	}
	if len(containers) == 0 {
		return fmt.Errorf("no se encontraron contenedores para el stack %q", projectName)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []string
	timeout := 15

	for _, c := range containers {
		wg.Add(1)
		go func(cid string) {
			defer wg.Done()
			if err := r.client.cli.ContainerRestart(ctx, cid, container.StopOptions{Timeout: &timeout}); err != nil {
				mu.Lock()
				short := cid
				if len(short) > 12 {
					short = short[:12]
				}
				errs = append(errs, fmt.Sprintf("%s: %v", short, err))
				mu.Unlock()
			}
		}(c.ID)
	}
	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("errores al reiniciar el stack: %s", strings.Join(errs, "; "))
	}
	return nil
}
