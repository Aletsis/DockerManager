package docker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"dockermanager/internal/domain"
	networkdomain "dockermanager/internal/domain/network"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	dockernetwork "github.com/docker/docker/api/types/network"
)

// NetworkRepository implements networkdomain.Repository using Docker SDK
type NetworkRepository struct {
	client *Client
}

// NewNetworkRepository creates a new Docker NetworkRepository
func NewNetworkRepository(client *Client) *NetworkRepository {
	return &NetworkRepository{client: client}
}

// List returns all docker networks enriched with connected containers and IPAM info
func (r *NetworkRepository) List(ctx context.Context) ([]networkdomain.Network, error) {
	nets, err := r.client.cli.NetworkList(ctx, dockernetwork.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error al listar redes Docker: %w", err)
	}

	// 1. Fetch all containers to discover network attachments (since NetworkList omits Containers map)
	containersList, _ := r.client.cli.ContainerList(ctx, container.ListOptions{All: true})
	netContainersMap := make(map[string]map[string]networkdomain.NetworkContainerRef)

	for _, c := range containersList {
		cName := ""
		if len(c.Names) > 0 {
			cName = strings.TrimPrefix(c.Names[0], "/")
		}
		shortCID := c.ID
		if len(shortCID) > 12 {
			shortCID = shortCID[:12]
		}

		if c.NetworkSettings != nil && c.NetworkSettings.Networks != nil {
			for netName, netEndpoint := range c.NetworkSettings.Networks {
				if netEndpoint == nil {
					continue
				}
				ref := networkdomain.NetworkContainerRef{
					ID:          shortCID,
					Name:        cName,
					State:       c.State,
					IPv4Address: netEndpoint.IPAddress,
					IPv6Address: netEndpoint.GlobalIPv6Address,
					MacAddress:  netEndpoint.MacAddress,
					EndpointID:  netEndpoint.EndpointID,
				}

				// Associate by Network ID
				if netEndpoint.NetworkID != "" {
					if netContainersMap[netEndpoint.NetworkID] == nil {
						netContainersMap[netEndpoint.NetworkID] = make(map[string]networkdomain.NetworkContainerRef)
					}
					netContainersMap[netEndpoint.NetworkID][c.ID] = ref
				}
				// Also associate by Network Name
				if netName != "" {
					if netContainersMap[netName] == nil {
						netContainersMap[netName] = make(map[string]networkdomain.NetworkContainerRef)
					}
					netContainersMap[netName][c.ID] = ref
				}
			}
		}
	}

	result := make([]networkdomain.Network, 0, len(nets))
	for _, net := range nets {
		shortID := net.ID
		if len(shortID) > 12 {
			shortID = shortID[:12]
		}

		isDefault := net.Name == "bridge" || net.Name == "host" || net.Name == "none"

		ipams := make([]networkdomain.IPAMConfig, 0, len(net.IPAM.Config))
		for _, cfg := range net.IPAM.Config {
			ipams = append(ipams, networkdomain.IPAMConfig{
				Driver:  net.IPAM.Driver,
				Subnet:  cfg.Subnet,
				Gateway: cfg.Gateway,
				IPRange: cfg.IPRange,
			})
		}

		// Merge containers discovered from ContainerList and net.Containers (if populated by daemon)
		containerSet := make(map[string]networkdomain.NetworkContainerRef)

		// A. From netContainersMap by network ID
		if byID, ok := netContainersMap[net.ID]; ok {
			for cid, ref := range byID {
				containerSet[cid] = ref
			}
		}
		// B. From netContainersMap by network Name
		if byName, ok := netContainersMap[net.Name]; ok {
			for cid, ref := range byName {
				containerSet[cid] = ref
			}
		}
		// C. From net.Containers if daemon provided them
		for cid, endpoint := range net.Containers {
			shortCID := cid
			if len(shortCID) > 12 {
				shortCID = shortCID[:12]
			}
			if _, exists := containerSet[cid]; !exists {
				containerSet[cid] = networkdomain.NetworkContainerRef{
					ID:          shortCID,
					Name:        endpoint.Name,
					IPv4Address: endpoint.IPv4Address,
					IPv6Address: endpoint.IPv6Address,
					MacAddress:  endpoint.MacAddress,
					EndpointID:  endpoint.EndpointID,
				}
			}
		}

		containers := make([]networkdomain.NetworkContainerRef, 0, len(containerSet))
		for _, ref := range containerSet {
			containers = append(containers, ref)
		}

		sort.Slice(containers, func(i, j int) bool {
			return strings.ToLower(containers[i].Name) < strings.ToLower(containers[j].Name)
		})

		createdStr := ""
		if !net.Created.IsZero() {
			createdStr = net.Created.Format(time.RFC3339)
		}

		result = append(result, networkdomain.Network{
			ID:              net.ID,
			ShortID:         shortID,
			Name:            net.Name,
			Driver:          net.Driver,
			Scope:           net.Scope,
			Internal:        net.Internal,
			Attachable:      net.Attachable,
			EnableIPv6:      net.EnableIPv6,
			IPAM:            ipams,
			Containers:      containers,
			ContainersCount: len(containers),
			Labels:          net.Labels,
			Options:         net.Options,
			Created:         createdStr,
			IsDefault:       isDefault,
		})
	}

	// Sort: user networks first (alphabetically), default networks at the end
	sort.Slice(result, func(i, j int) bool {
		if result[i].IsDefault && !result[j].IsDefault {
			return false
		}
		if !result[i].IsDefault && result[j].IsDefault {
			return true
		}
		return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
	})

	return result, nil
}

// Inspect returns formatted JSON of a docker network
func (r *NetworkRepository) Inspect(ctx context.Context, idOrName string) (string, error) {
	idOrName = strings.TrimSpace(idOrName)
	if idOrName == "" {
		return "", fmt.Errorf("el identificador de la red no puede estar vacío")
	}

	_, raw, err := r.client.cli.NetworkInspectWithRaw(ctx, idOrName, dockernetwork.InspectOptions{})
	if err != nil {
		return "", fmt.Errorf("error al inspeccionar red %q: %w", idOrName, err)
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, raw, "", "  "); err != nil {
		return string(raw), nil
	}
	return prettyJSON.String(), nil
}

// Create creates a new Docker network
func (r *NetworkRepository) Create(ctx context.Context, spec networkdomain.CreateNetworkSpec) (string, error) {
	name := strings.TrimSpace(spec.Name)
	if name == "" {
		return "", fmt.Errorf("el nombre de la red no puede estar vacío")
	}

	driver := strings.TrimSpace(spec.Driver)
	if driver == "" {
		driver = "bridge"
	}

	opts := dockernetwork.CreateOptions{
		Driver:     driver,
		Internal:   spec.Internal,
		Attachable: spec.Attachable,
		Labels:     spec.Labels,
		Options:    spec.Options,
	}

	if spec.EnableIPv6 {
		enableIPv6 := true
		opts.EnableIPv6 = &enableIPv6
	}

	if strings.TrimSpace(spec.Subnet) != "" || strings.TrimSpace(spec.Gateway) != "" || strings.TrimSpace(spec.IPRange) != "" {
		opts.IPAM = &dockernetwork.IPAM{
			Config: []dockernetwork.IPAMConfig{
				{
					Subnet:  strings.TrimSpace(spec.Subnet),
					Gateway: strings.TrimSpace(spec.Gateway),
					IPRange: strings.TrimSpace(spec.IPRange),
				},
			},
		}
	}

	resp, err := r.client.cli.NetworkCreate(ctx, name, opts)
	if err != nil {
		return "", fmt.Errorf("error al crear red %q: %w", name, err)
	}
	return resp.ID, nil
}

// Remove deletes a Docker network
func (r *NetworkRepository) Remove(ctx context.Context, idOrName string) error {
	idOrName = strings.TrimSpace(idOrName)
	if idOrName == "" {
		return fmt.Errorf("el identificador de la red no puede estar vacío")
	}

	if idOrName == "bridge" || idOrName == "host" || idOrName == "none" {
		return domain.ErrCannotRemoveDefaultNetwork
	}

	if err := r.client.cli.NetworkRemove(ctx, idOrName); err != nil {
		return fmt.Errorf("error al eliminar red %q: %w", idOrName, err)
	}
	return nil
}

// Prune removes unused Docker networks
func (r *NetworkRepository) Prune(ctx context.Context) (*networkdomain.PruneResult, error) {
	report, err := r.client.cli.NetworksPrune(ctx, filters.NewArgs())
	if err != nil {
		return nil, fmt.Errorf("error al limpiar redes inactivas: %w", err)
	}
	return &networkdomain.PruneResult{
		NetworksDeleted: report.NetworksDeleted,
	}, nil
}

// ConnectContainer attaches a container to a network
func (r *NetworkRepository) ConnectContainer(ctx context.Context, networkID, containerID string, ipAddress string) error {
	networkID = strings.TrimSpace(networkID)
	containerID = strings.TrimSpace(containerID)
	if networkID == "" || containerID == "" {
		return fmt.Errorf("networkID y containerID no pueden estar vacíos")
	}

	var endpointConfig *dockernetwork.EndpointSettings
	if strings.TrimSpace(ipAddress) != "" {
		endpointConfig = &dockernetwork.EndpointSettings{
			IPAMConfig: &dockernetwork.EndpointIPAMConfig{
				IPv4Address: strings.TrimSpace(ipAddress),
			},
		}
	}

	if err := r.client.cli.NetworkConnect(ctx, networkID, containerID, endpointConfig); err != nil {
		return fmt.Errorf("error al conectar contenedor a la red: %w", err)
	}
	return nil
}

// DisconnectContainer detaches a container from a network
func (r *NetworkRepository) DisconnectContainer(ctx context.Context, networkID, containerID string, force bool) error {
	networkID = strings.TrimSpace(networkID)
	containerID = strings.TrimSpace(containerID)
	if networkID == "" || containerID == "" {
		return fmt.Errorf("networkID y containerID no pueden estar vacíos")
	}

	if err := r.client.cli.NetworkDisconnect(ctx, networkID, containerID, force); err != nil {
		return fmt.Errorf("error al desconectar contenedor de la red: %w", err)
	}
	return nil
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
