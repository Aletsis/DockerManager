package network_test

import (
	"testing"

	"dockermanager/internal/domain/container"
	"dockermanager/internal/domain/network"
)

func TestNetworkGroupAggregate(t *testing.T) {
	group := network.Group{
		Name:      "custom-net",
		NetworkID: "net123",
		IsDefault: false,
		Containers: []container.Container{
			{ID: "c1", Name: "app1", State: container.StateRunning},
			{ID: "c2", Name: "app2", State: container.StatePaused},
			{ID: "c3", Name: "app3", State: container.StateExited},
		},
	}

	group.RecalculateCounts()

	if group.TotalCount != 3 {
		t.Errorf("expected TotalCount = 3, got %d", group.TotalCount)
	}
	if group.RunningCount != 1 {
		t.Errorf("expected RunningCount = 1, got %d", group.RunningCount)
	}

	// Change app2 from Paused to Running
	group.Containers[1].State = container.StateRunning
	group.RecalculateCounts()

	if group.RunningCount != 2 {
		t.Errorf("expected RunningCount = 2, got %d", group.RunningCount)
	}
}

func TestNetworkCanRemove(t *testing.T) {
	// Default network cannot be removed
	defaultNet := network.Network{
		Name:      "bridge",
		IsDefault: true,
	}
	if err := defaultNet.CanRemove(); err == nil {
		t.Errorf("expected error removing default network, got nil")
	}

	// Network with active containers cannot be removed
	inUseNet := network.Network{
		Name:      "my-net",
		IsDefault: false,
		Containers: []network.NetworkContainerRef{
			{ID: "c1", Name: "app1"},
		},
	}
	if !inUseNet.IsInUse() {
		t.Errorf("expected network to be in use")
	}
	if err := inUseNet.CanRemove(); err == nil {
		t.Errorf("expected error removing network in use, got nil")
	}

	// Unused custom network can be removed
	freeNet := network.Network{
		Name:       "free-net",
		IsDefault:  false,
		Containers: []network.NetworkContainerRef{},
	}
	if freeNet.IsInUse() {
		t.Errorf("expected network to not be in use")
	}
	if err := freeNet.CanRemove(); err != nil {
		t.Errorf("expected no error removing free network, got %v", err)
	}
}

func TestCreateNetworkSpecValidate(t *testing.T) {
	specInvalid := network.CreateNetworkSpec{
		Name:   "   ",
		Driver: "bridge",
	}
	if err := specInvalid.Validate(); err == nil {
		t.Errorf("expected error on empty network name, got nil")
	}

	specValid := network.CreateNetworkSpec{
		Name:   "custom_net",
		Driver: "bridge",
	}
	if err := specValid.Validate(); err != nil {
		t.Errorf("expected valid spec, got %v", err)
	}
}
