package stack_test

import (
	"context"
	"testing"

	stackapp "dockermanager/internal/application/stack"
	stackdomain "dockermanager/internal/domain/stack"
)

type mockStackRepo struct {
	lastAction string
	lastStack  string
	err        error
}

func (m *mockStackRepo) StartStack(ctx context.Context, projectName string) error {
	m.lastAction = "start"
	m.lastStack = projectName
	return m.err
}

func (m *mockStackRepo) StopStack(ctx context.Context, projectName string) error {
	m.lastAction = "stop"
	m.lastStack = projectName
	return m.err
}

func (m *mockStackRepo) RestartStack(ctx context.Context, projectName string) error {
	m.lastAction = "restart"
	m.lastStack = projectName
	return m.err
}

func (m *mockStackRepo) UpStack(ctx context.Context, req stackdomain.ComposeDeployRequest, onLog func(string)) error {
	m.lastAction = "up"
	m.lastStack = req.ProjectName
	if onLog != nil {
		onLog("Deploying test")
	}
	return m.err
}

func (m *mockStackRepo) DownStack(ctx context.Context, req stackdomain.ComposeDownRequest, onLog func(string)) error {
	m.lastAction = "down"
	m.lastStack = req.ProjectName
	if onLog != nil {
		onLog("Tearing down test")
	}
	return m.err
}

func (m *mockStackRepo) GetComposeFile(ctx context.Context, projectName, workingDir, configFile string) (*stackdomain.ComposeFileInfo, error) {
	m.lastAction = "get"
	m.lastStack = projectName
	return &stackdomain.ComposeFileInfo{
		ProjectName:  projectName,
		WorkingDir:   workingDir,
		ConfigFile:   configFile,
		Content:      "version: '3.8'",
		ExistsOnDisk: true,
	}, m.err
}

func (m *mockStackRepo) SaveComposeFile(ctx context.Context, workingDir, configFile, content string) (string, error) {
	m.lastAction = "save"
	return "/tmp/compose.yaml", m.err
}

func (m *mockStackRepo) GetDefaultStackDirectory(projectName string) (string, error) {
	m.lastAction = "getDefaultDir"
	return "/home/user/.dockermanager/stacks/" + projectName, m.err
}

func TestManageStackUseCase(t *testing.T) {
	mock := &mockStackRepo{}
	uc := stackapp.NewManageStackUseCase(mock)
	ctx := context.Background()

	_ = uc.StartStack(ctx, "app_stack")
	if mock.lastAction != "start" || mock.lastStack != "app_stack" {
		t.Fatalf("expected start app_stack, got %s %s", mock.lastAction, mock.lastStack)
	}

	_ = uc.StopStack(ctx, "app_stack")
	if mock.lastAction != "stop" || mock.lastStack != "app_stack" {
		t.Fatalf("expected stop app_stack, got %s %s", mock.lastAction, mock.lastStack)
	}

	_ = uc.RestartStack(ctx, "app_stack")
	if mock.lastAction != "restart" || mock.lastStack != "app_stack" {
		t.Fatalf("expected restart app_stack, got %s %s", mock.lastAction, mock.lastStack)
	}

	var logged []string
	err := uc.UpStack(ctx, stackdomain.ComposeDeployRequest{ProjectName: "demo_stack"}, func(l string) {
		logged = append(logged, l)
	})
	if err != nil || mock.lastAction != "up" || mock.lastStack != "demo_stack" || len(logged) == 0 {
		t.Fatalf("UpStack failed, got action: %s, stack: %s, logs: %v", mock.lastAction, mock.lastStack, logged)
	}

	err = uc.DownStack(ctx, stackdomain.ComposeDownRequest{ProjectName: "demo_stack"}, nil)
	if err != nil || mock.lastAction != "down" || mock.lastStack != "demo_stack" {
		t.Fatalf("DownStack failed, got action: %s, stack: %s", mock.lastAction, mock.lastStack)
	}

	info, err := uc.GetComposeFile(ctx, "demo_stack", "", "")
	if err != nil || mock.lastAction != "get" || info.Content == "" {
		t.Fatalf("GetComposeFile failed: %v", err)
	}

	saved, err := uc.SaveComposeFile(ctx, "/tmp", "", "content")
	if err != nil || mock.lastAction != "save" || saved == "" {
		t.Fatalf("SaveComposeFile failed: %v", err)
	}

	dir, err := uc.GetDefaultStackDirectory("demo_stack")
	if err != nil || mock.lastAction != "getDefaultDir" || dir == "" {
		t.Fatalf("GetDefaultStackDirectory failed: %v", err)
	}
}
