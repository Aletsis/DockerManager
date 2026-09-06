package container_test

import (
	"context"
	"errors"
	"testing"

	containerapp "dockermanager/internal/application/container"
	"dockermanager/internal/domain"
	containerdomain "dockermanager/internal/domain/container"
)

// mockContainerRepo implements containerdomain.Repository for testing
type mockContainerRepo struct {
	containers []containerdomain.Container
	stats      *containerdomain.Stats
	inspect    string
	logs       string
	createResp *containerdomain.CreateResult
	err        error

	lastAction   string
	lastID       string
	lastTail     int
	lastForce    bool
	lastCreated  containerdomain.CreateSpec
}

func (m *mockContainerRepo) List(ctx context.Context, all bool) ([]containerdomain.Container, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.containers, nil
}

func (m *mockContainerRepo) FindByID(ctx context.Context, id string) (*containerdomain.Container, error) {
	for _, c := range m.containers {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, domain.ErrContainerNotFound
}

func (m *mockContainerRepo) Start(ctx context.Context, id string) error {
	m.lastAction = "start"
	m.lastID = id
	return m.err
}

func (m *mockContainerRepo) Stop(ctx context.Context, id string) error {
	m.lastAction = "stop"
	m.lastID = id
	return m.err
}

func (m *mockContainerRepo) Restart(ctx context.Context, id string) error {
	m.lastAction = "restart"
	m.lastID = id
	return m.err
}

func (m *mockContainerRepo) Pause(ctx context.Context, id string) error {
	m.lastAction = "pause"
	m.lastID = id
	return m.err
}

func (m *mockContainerRepo) Unpause(ctx context.Context, id string) error {
	m.lastAction = "unpause"
	m.lastID = id
	return m.err
}

func (m *mockContainerRepo) Remove(ctx context.Context, id string, force bool) error {
	m.lastAction = "remove"
	m.lastID = id
	m.lastForce = force
	return m.err
}

func (m *mockContainerRepo) GetLogs(ctx context.Context, id string, tail int) (string, error) {
	m.lastID = id
	m.lastTail = tail
	return m.logs, m.err
}

func (m *mockContainerRepo) GetStats(ctx context.Context, id string) (*containerdomain.Stats, error) {
	m.lastID = id
	return m.stats, m.err
}

func (m *mockContainerRepo) Inspect(ctx context.Context, id string) (string, error) {
	m.lastID = id
	return m.inspect, m.err
}

func (m *mockContainerRepo) Create(ctx context.Context, spec containerdomain.CreateSpec) (*containerdomain.CreateResult, error) {
	m.lastCreated = spec
	return m.createResp, m.err
}

func TestListContainersUseCase(t *testing.T) {
	mock := &mockContainerRepo{
		containers: []containerdomain.Container{
			{ID: "c1", Name: "app1"},
			{ID: "c2", Name: "app2"},
		},
	}
	uc := containerapp.NewListContainersUseCase(mock)

	res, err := uc.Execute(context.Background(), true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 2 {
		t.Errorf("expected 2 containers, got %d", len(res))
	}
}

func TestManageLifecycleUseCase(t *testing.T) {
	mock := &mockContainerRepo{}
	uc := containerapp.NewManageLifecycleUseCase(mock)
	ctx := context.Background()

	_ = uc.Start(ctx, "c1")
	if mock.lastAction != "start" || mock.lastID != "c1" {
		t.Errorf("expected start on c1, got %s on %s", mock.lastAction, mock.lastID)
	}

	_ = uc.Stop(ctx, "c2")
	if mock.lastAction != "stop" || mock.lastID != "c2" {
		t.Errorf("expected stop on c2, got %s on %s", mock.lastAction, mock.lastID)
	}

	_ = uc.Restart(ctx, "c3")
	if mock.lastAction != "restart" || mock.lastID != "c3" {
		t.Errorf("expected restart on c3, got %s on %s", mock.lastAction, mock.lastID)
	}

	_ = uc.Pause(ctx, "c4")
	if mock.lastAction != "pause" || mock.lastID != "c4" {
		t.Errorf("expected pause on c4, got %s on %s", mock.lastAction, mock.lastID)
	}

	_ = uc.Unpause(ctx, "c5")
	if mock.lastAction != "unpause" || mock.lastID != "c5" {
		t.Errorf("expected unpause on c5, got %s on %s", mock.lastAction, mock.lastID)
	}

	_ = uc.Remove(ctx, "c6", true)
	if mock.lastAction != "remove" || mock.lastID != "c6" || !mock.lastForce {
		t.Errorf("expected remove on c6 with force, got %s on %s force=%v", mock.lastAction, mock.lastID, mock.lastForce)
	}
}

func TestCreateContainerUseCase(t *testing.T) {
	mock := &mockContainerRepo{
		createResp: &containerdomain.CreateResult{ID: "new_c1"},
	}
	uc := containerapp.NewCreateContainerUseCase(mock)
	ctx := context.Background()

	// Validation failure (empty image)
	_, err := uc.Execute(ctx, containerdomain.CreateSpec{Image: ""})
	if !errors.Is(err, domain.ErrInvalidContainerSpec) {
		t.Errorf("expected ErrInvalidContainerSpec, got %v", err)
	}

	// Success
	res, err := uc.Execute(ctx, containerdomain.CreateSpec{Image: "redis:alpine", Name: "my_redis"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != "new_c1" {
		t.Errorf("expected new_c1, got %s", res.ID)
	}
	if mock.lastCreated.Name != "my_redis" {
		t.Errorf("expected spec passed to repo with name my_redis, got %s", mock.lastCreated.Name)
	}
}

func TestGetTelemetryUseCase(t *testing.T) {
	mock := &mockContainerRepo{
		logs:    "system ready\nlistening on 80\n",
		stats:   &containerdomain.Stats{ID: "c1", CPUPercentage: 12.5},
		inspect: "{\"Id\": \"c1\"}",
	}
	uc := containerapp.NewGetTelemetryUseCase(mock)
	ctx := context.Background()

	logs, err := uc.GetLogs(ctx, "c1", 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if logs != "system ready\nlistening on 80\n" {
		t.Errorf("unexpected logs: %s", logs)
	}
	if mock.lastTail != 100 {
		t.Errorf("expected tail 100, got %d", mock.lastTail)
	}

	stats, err := uc.GetStats(ctx, "c1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.CPUPercentage != 12.5 {
		t.Errorf("expected CPU 12.5, got %f", stats.CPUPercentage)
	}

	insp, err := uc.Inspect(ctx, "c1")
	if err != nil {
		t.Fatalf("unexpected error inspecting: %v", err)
	}
	if insp != "{\"Id\": \"c1\"}" {
		t.Errorf("expected {\"Id\": \"c1\"}, got %s", insp)
	}
}
