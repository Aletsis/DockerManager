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
