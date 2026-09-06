package container

import (
	"dockermanager/internal/domain"
	"strings"
)

// Container is the rich domain entity representing a Docker/OCI container
type Container struct {
	ID                string            `json:"id"`
	ShortID           string            `json:"shortId"`
	Names             []string          `json:"names"`
	Name              string            `json:"name"`
	Image             string            `json:"image"`
	ImageID           string            `json:"imageId"`
	Command           string            `json:"command"`
	Created           int64             `json:"created"`
	State             State             `json:"state"`
	Status            string            `json:"status"`
	Ports             []PortMapping     `json:"ports"`
	Networks          []NetworkInfo     `json:"networks,omitempty"`
	SizeRw            int64             `json:"sizeRw"`
	SizeRootFs        int64             `json:"sizeRootFs"`
	Labels            map[string]string `json:"labels,omitempty"`
	ComposeProject    string            `json:"composeProject,omitempty"`
	ComposeService    string            `json:"composeService,omitempty"`
	ComposeWorkingDir string            `json:"composeWorkingDir,omitempty"`
	ComposeConfigFile string            `json:"composeConfigFile,omitempty"`
}

// IsRunning returns true if the container is actively executing
func (c *Container) IsRunning() bool {
	return c.State == StateRunning
}

// IsPaused returns true if the container execution is frozen
func (c *Container) IsPaused() bool {
	return c.State == StatePaused
}

// CanStart checks domain invariants before attempting to start
func (c *Container) CanStart() error {
	if c.IsRunning() {
		return domain.ErrContainerAlreadyRunning
	}
	if c.IsPaused() {
		return domain.ErrContainerPaused
	}
	return nil
}

// CanStop checks domain invariants before attempting to stop
func (c *Container) CanStop() error {
	if !c.IsRunning() && !c.IsPaused() {
		return domain.ErrContainerNotRunning
	}
	return nil
}

// CanPause checks domain invariants before attempting to pause
func (c *Container) CanPause() error {
	if !c.IsRunning() {
		return domain.ErrContainerNotRunning
	}
	return nil
}

// CanUnpause checks domain invariants before attempting to unpause
func (c *Container) CanUnpause() error {
	if !c.IsPaused() {
		return domain.ErrContainerNotPaused
	}
	return nil
}

// CanRemove checks whether the container can be safely deleted
func (c *Container) CanRemove(force bool) error {
	if c.IsRunning() && !force {
		return domain.ErrCannotRemoveRunning
	}
	return nil
}

// BelongsToCompose returns true if the container is part of a Docker Compose project
func (c *Container) BelongsToCompose() bool {
	return c.ComposeProject != ""
}

// BelongsToNetwork returns true if the container is attached to the given network name
func (c *Container) BelongsToNetwork(networkName string) bool {
	for _, net := range c.Networks {
		if net.NetworkName == networkName {
			return true
		}
	}
	return false
}

// ValidateCreateSpec validates input specifications for container creation
func ValidateCreateSpec(spec CreateSpec) error {
	if strings.TrimSpace(spec.Image) == "" {
		return domain.ErrInvalidContainerSpec
	}
	return nil
}
