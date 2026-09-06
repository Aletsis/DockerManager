package docker_test

import (
	"os"
	"path/filepath"
	"testing"

	"dockermanager/internal/infrastructure/docker"
)

func TestResolveComposeFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dm_compose_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 1. Test when no files exist
	_, found := docker.ResolveComposeFile(tempDir, "")
	if found {
		t.Fatalf("expected not found in empty dir")
	}

	// 2. Create docker-compose.yml
	composePath := filepath.Join(tempDir, "docker-compose.yml")
	if err := os.WriteFile(composePath, []byte("services: {}"), 0644); err != nil {
		t.Fatalf("failed to write compose file: %v", err)
	}

	resolved, found := docker.ResolveComposeFile(tempDir, "")
	if !found || resolved != composePath {
		t.Fatalf("expected resolved path %s, got %s (found: %v)", composePath, resolved, found)
	}

	// 3. Test with comma-separated config_files label
	overridePath := filepath.Join(tempDir, "docker-compose.override.yml")
	_ = os.WriteFile(overridePath, []byte("services: {}"), 0644)

	commaSeparated := composePath + "," + overridePath
	resolved2, found2 := docker.ResolveComposeFile(tempDir, commaSeparated)
	if !found2 || resolved2 != composePath {
		t.Fatalf("expected comma-separated first match %s, got %s", composePath, resolved2)
	}
}

func TestDiscoverDockerfiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "dm_dockerfile_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 1. Empty dir should return empty
	dfs := docker.DiscoverDockerfiles(tempDir, "", nil)
	if len(dfs) != 0 {
		t.Fatalf("expected 0 dockerfiles, got %d", len(dfs))
	}

	// 2. Root Dockerfile
	rootDf := filepath.Join(tempDir, "Dockerfile")
	_ = os.WriteFile(rootDf, []byte("FROM alpine:latest\n"), 0644)

	// 3. Service Dockerfile
	svcDir := filepath.Join(tempDir, "backend")
	_ = os.MkdirAll(svcDir, 0755)
	svcDf := filepath.Join(svcDir, "Dockerfile")
	_ = os.WriteFile(svcDf, []byte("FROM golang:alpine\n"), 0644)

	// 4. Ignored node_modules Dockerfile
	nodeModulesDir := filepath.Join(tempDir, "node_modules", "some-pkg")
	_ = os.MkdirAll(nodeModulesDir, 0755)
	_ = os.WriteFile(filepath.Join(nodeModulesDir, "Dockerfile"), []byte("FROM node\n"), 0644)

	discovered := docker.DiscoverDockerfiles(tempDir, "", []string{"backend"})
	if len(discovered) != 2 {
		t.Fatalf("expected 2 discovered dockerfiles (root + backend, ignoring node_modules), got %d: %+v", len(discovered), discovered)
	}

	foundRoot := false
	foundBackend := false
	for _, df := range discovered {
		if df.Name == "Dockerfile" && df.Content == "FROM alpine:latest\n" {
			foundRoot = true
		}
		if df.Name == "backend/Dockerfile" && df.ServiceName == "backend" {
			foundBackend = true
		}
	}

	if !foundRoot || !foundBackend {
		t.Fatalf("missing expected Dockerfile: root=%v, backend=%v (list: %+v)", foundRoot, foundBackend, discovered)
	}
}
