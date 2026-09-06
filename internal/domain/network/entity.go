package network

import (
	"dockermanager/internal/domain"
	"dockermanager/internal/domain/container"
	"strings"
)

// NetworkContainerRef represents a container connected to a network
type NetworkContainerRef struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	State       string `json:"state,omitempty"`
	IPv4Address string `json:"ipv4Address"`
	IPv6Address string `json:"ipv6Address"`
	MacAddress  string `json:"macAddress"`
	EndpointID  string `json:"endpointId"`
}

// IPAMConfig represents IP Address Management configuration of a network
type IPAMConfig struct {
	Driver  string `json:"driver,omitempty"`
	Subnet  string `json:"subnet,omitempty"`
	Gateway string `json:"gateway,omitempty"`
	IPRange string `json:"ipRange,omitempty"`
}

// Network represents a Docker network
type Network struct {
	ID              string                `json:"id"`
	ShortID         string                `json:"shortId"`
	Name            string                `json:"name"`
	Driver          string                `json:"driver"`
	Scope           string                `json:"scope"`
	Internal        bool                  `json:"internal"`
	Attachable      bool                  `json:"attachable"`
	EnableIPv6      bool                  `json:"enableIPv6"`
	IPAM            []IPAMConfig          `json:"ipam"`
	Containers      []NetworkContainerRef `json:"containers"`
	ContainersCount int                   `json:"containersCount"`
	Labels          map[string]string     `json:"labels"`
	Options         map[string]string     `json:"options"`
	Created         string                `json:"created"`
	IsDefault       bool                  `json:"isDefault"`
}

// IsInUse returns true if at least one container is currently attached
func (n *Network) IsInUse() bool {
	return len(n.Containers) > 0
}

// CanRemove checks whether a network can be safely deleted
func (n *Network) CanRemove() error {
	if n.IsDefault {
		return domain.ErrCannotRemoveDefaultNetwork
	}
	if n.IsInUse() {
		return domain.ErrNetworkInUse
	}
	return nil
}

// CreateNetworkSpec defines parameters for creating a new Docker network
type CreateNetworkSpec struct {
	Name       string            `json:"name"`
	Driver     string            `json:"driver"`
	Subnet     string            `json:"subnet,omitempty"`
	Gateway    string            `json:"gateway,omitempty"`
	IPRange    string            `json:"ipRange,omitempty"`
	Internal   bool              `json:"internal"`
	Attachable bool              `json:"attachable"`
	EnableIPv6 bool              `json:"enableIPv6"`
	Labels     map[string]string `json:"labels,omitempty"`
	Options    map[string]string `json:"options,omitempty"`
}

// Validate checks for mandatory fields in CreateNetworkSpec
func (s *CreateNetworkSpec) Validate() error {
	if strings.TrimSpace(s.Name) == "" {
		return domain.ErrNetworkNotFound
	}
	return nil
}

// PruneResult contains the result of network pruning
type PruneResult struct {
	NetworksDeleted []string `json:"networksDeleted"`
}

// Group represents a Docker network grouping connected containers
type Group struct {
	Name         string                `json:"name"`
	NetworkID    string                `json:"networkId,omitempty"`
	IsDefault    bool                  `json:"isDefault,omitempty"`
	Containers   []container.Container `json:"containers"`
	RunningCount int                   `json:"runningCount"`
	TotalCount   int                   `json:"totalCount"`
}

// RecalculateCounts updates running and total container tallies
func (g *Group) RecalculateCounts() {
	g.TotalCount = len(g.Containers)
	g.RunningCount = 0
	for _, c := range g.Containers {
		if c.IsRunning() {
			g.RunningCount++
		}
	}
}
