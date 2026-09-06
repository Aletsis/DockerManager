package network_test

import (
	"context"
	"testing"

	networkapp "dockermanager/internal/application/network"
)

type mockNetworkRepo struct {
	lastAction  string
	lastNetwork string
	err         error
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
}
