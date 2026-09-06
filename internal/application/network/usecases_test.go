package network_test

import (
	"context"
	"testing"

	networkapp "dockermanager/internal/application/network"
	networkdomain "dockermanager/internal/domain/network"
)

type mockNetworkRepo struct {
	lastAction    string
	lastNetwork   string
	lastContainer string
	lastIP        string
	lastForce     bool
	specToCreate  networkdomain.CreateNetworkSpec
	err           error
}

func (m *mockNetworkRepo) List(ctx context.Context) ([]networkdomain.Network, error) {
	m.lastAction = "list"
	return []networkdomain.Network{{ID: "net1", Name: "bridge"}}, m.err
}

func (m *mockNetworkRepo) Inspect(ctx context.Context, idOrName string) (string, error) {
	m.lastAction = "inspect"
	m.lastNetwork = idOrName
	return `{"Id": "` + idOrName + `"}`, m.err
}

func (m *mockNetworkRepo) Create(ctx context.Context, spec networkdomain.CreateNetworkSpec) (string, error) {
	m.lastAction = "create"
	m.specToCreate = spec
	return "net_created_123", m.err
}

func (m *mockNetworkRepo) Remove(ctx context.Context, idOrName string) error {
	m.lastAction = "remove"
	m.lastNetwork = idOrName
	return m.err
}

func (m *mockNetworkRepo) Prune(ctx context.Context) (*networkdomain.PruneResult, error) {
	m.lastAction = "prune"
	return &networkdomain.PruneResult{NetworksDeleted: []string{"old_net"}}, m.err
}

func (m *mockNetworkRepo) ConnectContainer(ctx context.Context, networkID, containerID string, ipAddress string) error {
	m.lastAction = "connect"
	m.lastNetwork = networkID
	m.lastContainer = containerID
	m.lastIP = ipAddress
	return m.err
}

func (m *mockNetworkRepo) DisconnectContainer(ctx context.Context, networkID, containerID string, force bool) error {
	m.lastAction = "disconnect"
	m.lastNetwork = networkID
	m.lastContainer = containerID
	m.lastForce = force
	return m.err
}

func (m *mockNetworkRepo) StartNetwork(ctx context.Context, networkName string) error {
	m.lastAction = "start"
	m.lastNetwork = networkName
	return m.err
}

func (m *mockNetworkRepo) StopNetwork(ctx context.Context, networkName string) error {
	m.lastAction = "stop"
	m.lastNetwork = networkName
	return m.err
}

func (m *mockNetworkRepo) RestartNetwork(ctx context.Context, networkName string) error {
	m.lastAction = "restart"
	m.lastNetwork = networkName
	return m.err
}

func TestManageNetworkUseCase(t *testing.T) {
	mock := &mockNetworkRepo{}
	uc := networkapp.NewManageNetworkUseCase(mock)
	ctx := context.Background()

	_ = uc.StartNetwork(ctx, "dev_net")
	if mock.lastAction != "start" || mock.lastNetwork != "dev_net" {
		t.Fatalf("expected start dev_net, got %s %s", mock.lastAction, mock.lastNetwork)
	}

	_ = uc.StopNetwork(ctx, "dev_net")
	if mock.lastAction != "stop" || mock.lastNetwork != "dev_net" {
		t.Fatalf("expected stop dev_net, got %s %s", mock.lastAction, mock.lastNetwork)
	}

	_ = uc.RestartNetwork(ctx, "dev_net")
	if mock.lastAction != "restart" || mock.lastNetwork != "dev_net" {
		t.Fatalf("expected restart dev_net, got %s %s", mock.lastAction, mock.lastNetwork)
	}

	_ = uc.Remove(ctx, "dev_net")
	if mock.lastAction != "remove" || mock.lastNetwork != "dev_net" {
		t.Fatalf("expected remove dev_net, got %s %s", mock.lastAction, mock.lastNetwork)
	}

	res, _ := uc.Prune(ctx)
	if mock.lastAction != "prune" || len(res.NetworksDeleted) != 1 {
		t.Fatalf("expected prune, got %s %v", mock.lastAction, res)
	}

	_ = uc.Connect(ctx, "dev_net", "cid1", "172.20.0.5")
	if mock.lastAction != "connect" || mock.lastNetwork != "dev_net" || mock.lastContainer != "cid1" || mock.lastIP != "172.20.0.5" {
		t.Fatalf("expected connect dev_net cid1, got %s %s %s", mock.lastAction, mock.lastNetwork, mock.lastContainer)
	}

	_ = uc.Disconnect(ctx, "dev_net", "cid1", true)
	if mock.lastAction != "disconnect" || mock.lastNetwork != "dev_net" || mock.lastContainer != "cid1" || !mock.lastForce {
		t.Fatalf("expected disconnect dev_net cid1, got %s %s %s", mock.lastAction, mock.lastNetwork, mock.lastContainer)
	}

	inspectRes, _ := uc.Inspect(ctx, "dev_net")
	if mock.lastAction != "inspect" || mock.lastNetwork != "dev_net" || inspectRes == "" {
		t.Fatalf("expected inspect dev_net, got %s %s", mock.lastAction, mock.lastNetwork)
	}
}

func TestListNetworksUseCase(t *testing.T) {
	mock := &mockNetworkRepo{}
	uc := networkapp.NewListNetworksUseCase(mock)
	ctx := context.Background()

	list, err := uc.Execute(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.lastAction != "list" || len(list) != 1 {
		t.Fatalf("expected list with 1 network, got %v", list)
	}
}

func TestCreateNetworkUseCase(t *testing.T) {
	mock := &mockNetworkRepo{}
	uc := networkapp.NewCreateNetworkUseCase(mock)
	ctx := context.Background()

	// Invalid spec (empty name)
	_, err := uc.Execute(ctx, networkdomain.CreateNetworkSpec{Name: ""})
	if err == nil {
		t.Fatalf("expected error on empty network name, got nil")
	}

	// Valid spec
	id, err := uc.Execute(ctx, networkdomain.CreateNetworkSpec{
		Name:   "custom_net",
		Driver: "bridge",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "net_created_123" || mock.lastAction != "create" {
		t.Fatalf("expected net_created_123, got %s", id)
	}
}
