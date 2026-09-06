package docker

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

// NetworkRepository implements networkdomain.Repository using Docker SDK
type NetworkRepository struct {
	client *Client
}

// NewNetworkRepository creates a new Docker NetworkRepository
func NewNetworkRepository(client *Client) *NetworkRepository {
	return &NetworkRepository{client: client}
}

// StartNetwork starts all stopped containers connected to a Docker network
func (r *NetworkRepository) StartNetwork(ctx context.Context, networkName string) error {
	networkName = strings.TrimSpace(networkName)
	if networkName == "" {
		return fmt.Errorf("el nombre de la red no puede estar vacío")
	}

	filterArgs := filters.NewArgs()
	filterArgs.Add("network", networkName)

	containers, err := r.client.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filterArgs})
	if err != nil {
		return fmt.Errorf("error al listar contenedores de la red %s: %w", networkName, err)
	}
	if len(containers) == 0 {
		return fmt.Errorf("no se encontraron contenedores para la red %q", networkName)
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
		return fmt.Errorf("errores al iniciar contenedores de la red: %s", strings.Join(errs, "; "))
	}
	return nil
}

// StopNetwork stops all running containers connected to a Docker network
func (r *NetworkRepository) StopNetwork(ctx context.Context, networkName string) error {
	networkName = strings.TrimSpace(networkName)
	if networkName == "" {
		return fmt.Errorf("el nombre de la red no puede estar vacío")
	}

	filterArgs := filters.NewArgs()
	filterArgs.Add("network", networkName)

	containers, err := r.client.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filterArgs})
	if err != nil {
		return fmt.Errorf("error al listar contenedores de la red %s: %w", networkName, err)
	}
	if len(containers) == 0 {
		return fmt.Errorf("no se encontraron contenedores para la red %q", networkName)
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
		return fmt.Errorf("errores al detener contenedores de la red: %s", strings.Join(errs, "; "))
	}
	return nil
}

// RestartNetwork restarts all containers connected to a Docker network
func (r *NetworkRepository) RestartNetwork(ctx context.Context, networkName string) error {
	networkName = strings.TrimSpace(networkName)
	if networkName == "" {
		return fmt.Errorf("el nombre de la red no puede estar vacío")
	}

	filterArgs := filters.NewArgs()
	filterArgs.Add("network", networkName)

	containers, err := r.client.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filterArgs})
	if err != nil {
		return fmt.Errorf("error al listar contenedores de la red %s: %w", networkName, err)
	}
	if len(containers) == 0 {
		return fmt.Errorf("no se encontraron contenedores para la red %q", networkName)
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
		return fmt.Errorf("errores al reiniciar contenedores de la red: %s", strings.Join(errs, "; "))
	}
	return nil
}
