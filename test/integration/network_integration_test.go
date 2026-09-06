package integration_test

import (
	"context"
	"testing"

	dockerinfra "dockermanager/internal/infrastructure/docker"
)

func TestNetworkRepositoryListWithContainers(t *testing.T) {
	client := getTestDockerClient(t)
	repo := dockerinfra.NewNetworkRepository(client)

	ctx := context.Background()
	nets, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("failed to list networks: %v", err)
	}

	if len(nets) == 0 {
		t.Fatalf("expected at least 1 network (bridge/host/none), got 0")
	}

	totalContainers := 0
	for _, n := range nets {
		totalContainers += n.ContainersCount
		t.Logf("Network %s (%s) has %d containers: %+v", n.Name, n.Driver, n.ContainersCount, n.Containers)
	}

	t.Logf("Total networks: %d, Total containers attached across networks: %d", len(nets), totalContainers)
}
