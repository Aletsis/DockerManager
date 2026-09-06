package docker

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

// StartNetwork starts all non-running containers connected to the specified Docker network
func (s *Service) StartNetwork(ctx context.Context, networkName string) error {
	networkName = strings.TrimSpace(networkName)
	if networkName == "" {
		return fmt.Errorf("el nombre de la red no puede estar vacío")
	}

	filterArgs := filters.NewArgs()
	filterArgs.Add("network", networkName)

	containers, err := s.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filterArgs})
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
			if err := s.StartContainer(ctx, cid); err != nil {
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

// StopNetwork stops all running containers connected to the specified Docker network
func (s *Service) StopNetwork(ctx context.Context, networkName string) error {
	networkName = strings.TrimSpace(networkName)
	if networkName == "" {
		return fmt.Errorf("el nombre de la red no puede estar vacío")
	}

	filterArgs := filters.NewArgs()
	filterArgs.Add("network", networkName)

	containers, err := s.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filterArgs})
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
		if c.State != "running" {
			continue
		}
		wg.Add(1)
		go func(cid string) {
			defer wg.Done()
			if err := s.StopContainer(ctx, cid); err != nil {
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

// RestartNetwork restarts all containers connected to the specified Docker network
func (s *Service) RestartNetwork(ctx context.Context, networkName string) error {
	networkName = strings.TrimSpace(networkName)
	if networkName == "" {
		return fmt.Errorf("el nombre de la red no puede estar vacío")
	}

	filterArgs := filters.NewArgs()
	filterArgs.Add("network", networkName)

	containers, err := s.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filterArgs})
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
		wg.Add(1)
		go func(cid string) {
			defer wg.Done()
			if err := s.RestartContainer(ctx, cid); err != nil {
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
