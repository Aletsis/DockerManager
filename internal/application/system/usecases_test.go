package system_test

import (
	"context"
	"testing"

	systemapp "dockermanager/internal/application/system"
	systemdomain "dockermanager/internal/domain/system"
)

type mockSystemRepo struct {
	overview *systemdomain.Overview
	err      error
}

func (m *mockSystemRepo) GetOverview(ctx context.Context) (*systemdomain.Overview, error) {
	return m.overview, m.err
}

func TestGetOverviewUseCase(t *testing.T) {
	mock := &mockSystemRepo{
		overview: &systemdomain.Overview{
			Containers:        5,
			ContainersRunning: 3,
			ServerVersion:     "24.0.5",
			OperatingSystem:   "Linux",
		},
	}
	uc := systemapp.NewGetOverviewUseCase(mock)
	res, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Containers != 5 || res.ContainersRunning != 3 || res.ServerVersion != "24.0.5" {
		t.Fatalf("unexpected overview data: %+v", res)
	}
}
