package docker

import (
	"testing"

	"github.com/docker/docker/api/types/container"
)

func TestParseRepoTag(t *testing.T) {
	tests := []struct {
		input       string
		wantRepo    string
		wantTag     string
		wantDangling bool
	}{
		{"nginx:alpine", "nginx", "alpine", false},
		{"localhost:5000/my-app:1.0.0", "localhost:5000/my-app", "1.0.0", false},
		{"redis", "redis", "latest", false},
		{"<none>:<none>", "<none>", "<none>", true},
		{"<none>", "<none>", "<none>", true},
		{"", "<none>", "<none>", true},
	}

	for _, tt := range tests {
		repo, tag, dangling := ParseRepoTag(tt.input)
		if repo != tt.wantRepo || tag != tt.wantTag || dangling != tt.wantDangling {
			t.Errorf("ParseRepoTag(%q) = (%q, %q, %v); want (%q, %q, %v)",
				tt.input, repo, tag, dangling, tt.wantRepo, tt.wantTag, tt.wantDangling)
		}
	}
}

func TestBuildContainerConfig(t *testing.T) {
	// 1. Error on empty image
	_, _, err := BuildContainerConfig(CreateContainerRequest{Image: "   "})
	if err == nil {
		t.Errorf("Expected error for empty image, got nil")
	}

	// 2. Valid request with ports, volumes, and env
	req := CreateContainerRequest{
		Image:         "nginx:alpine",
		Name:          "web_test",
		Ports:         []string{"8080:80/tcp", "9090:90/udp"},
		Volumes:       []string{"/tmp/data:/data:rw"},
		Env:           []string{"ENV_VAR=123"},
		RestartPolicy: "always",
	}

	cfg, hostCfg, err := BuildContainerConfig(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if cfg.Image != "nginx:alpine" {
		t.Errorf("Expected Image 'nginx:alpine', got %q", cfg.Image)
	}
	if len(cfg.Env) != 1 || cfg.Env[0] != "ENV_VAR=123" {
		t.Errorf("Unexpected Env: %v", cfg.Env)
	}
	if len(hostCfg.Binds) != 1 || hostCfg.Binds[0] != "/tmp/data:/data:rw" {
		t.Errorf("Unexpected Binds: %v", hostCfg.Binds)
	}
	if hostCfg.RestartPolicy.Name != "always" {
		t.Errorf("Unexpected RestartPolicy: %v", hostCfg.RestartPolicy.Name)
	}

	// 3. Invalid port format
	invalidReq := CreateContainerRequest{
		Image: "redis:alpine",
		Ports: []string{"invalid_port_spec"},
	}
	_, _, err = BuildContainerConfig(invalidReq)
	if err == nil {
		t.Errorf("Expected error for invalid port spec, got nil")
	}
}

func TestCalculateContainerStats(t *testing.T) {
	stats := &container.StatsResponse{
		ID:   "c123",
		Name: "my_container",
	}
	stats.CPUStats.CPUUsage.TotalUsage = 2000000
	stats.PreCPUStats.CPUUsage.TotalUsage = 1000000
	stats.CPUStats.SystemUsage = 10000000
	stats.PreCPUStats.SystemUsage = 5000000
	stats.CPUStats.OnlineCPUs = 2

	stats.MemoryStats.Usage = 104857600 // 100 MB
	stats.MemoryStats.Limit = 524288000 // 500 MB
	stats.MemoryStats.Stats = map[string]uint64{
		"cache": 4857600, // 100MB - ~4.85MB = 100000000 (roughly)
	}

	res := CalculateContainerStats("c123", "my_container", stats)
	if res.ID != "c123" || res.Name != "my_container" {
		t.Errorf("Unexpected ID or Name: %v, %v", res.ID, res.Name)
	}

	// CPU = (1,000,000 / 5,000,000) * 2 * 100 = 40.0%
	if res.CPUPercentage < 39.9 || res.CPUPercentage > 40.1 {
		t.Errorf("Expected CPU ~40.0, got %f", res.CPUPercentage)
	}

	// Mem Usage = 104857600 - 4857600 = 100000000
	if res.MemoryUsage != 100000000 {
		t.Errorf("Expected MemUsage 100000000, got %d", res.MemoryUsage)
	}

	// Mem Percentage = (100000000 / 524288000) * 100 ~= 19.07%
	if res.MemoryPercentage < 19.0 || res.MemoryPercentage > 19.2 {
		t.Errorf("Expected MemPercentage ~19.07, got %f", res.MemoryPercentage)
	}
}
