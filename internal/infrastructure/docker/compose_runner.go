package docker

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	stackdomain "dockermanager/internal/domain/stack"
)

// ComposeRunner manages executing Docker Compose CLI commands
type ComposeRunner struct {
	mu           sync.Mutex
	detectedCmd  []string
	hasEvaluated bool
}

// NewComposeRunner creates a new ComposeRunner
func NewComposeRunner() *ComposeRunner {
	return &ComposeRunner{}
}

// detectCommand checks if 'docker compose' or 'docker-compose' is available
func (r *ComposeRunner) detectCommand(ctx context.Context) ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.hasEvaluated && len(r.detectedCmd) > 0 {
		return r.detectedCmd, nil
	}

	// 1. Try 'docker compose version'
	cmdDockerCompose := exec.CommandContext(ctx, "docker", "compose", "version")
	if err := cmdDockerCompose.Run(); err == nil {
		r.detectedCmd = []string{"docker", "compose"}
		r.hasEvaluated = true
		return r.detectedCmd, nil
	}

	// 2. Fallback to standalone 'docker-compose version'
	cmdStandalone := exec.CommandContext(ctx, "docker-compose", "version")
	if err := cmdStandalone.Run(); err == nil {
		r.detectedCmd = []string{"docker-compose"}
		r.hasEvaluated = true
		return r.detectedCmd, nil
	}

	r.hasEvaluated = true
	return nil, errors.New("docker compose no está disponible en el sistema. Asegúrate de tener instalado el plugin 'docker compose' o el comando 'docker-compose'")
}

// Run executes a compose command with the specified arguments and streams output line by line
func (r *ComposeRunner) Run(ctx context.Context, workingDir string, args []string, onLog func(string)) error {
	baseCmd, err := r.detectCommand(ctx)
	if err != nil {
		return err
	}

	var fullArgs []string
	var bin string

	if len(baseCmd) == 1 {
		bin = baseCmd[0]
		fullArgs = args
	} else {
		bin = baseCmd[0]
		fullArgs = append(baseCmd[1:], args...)
	}

	cmd := exec.CommandContext(ctx, bin, fullArgs...)
	if workingDir != "" {
		if fi, statErr := os.Stat(workingDir); statErr == nil && fi.IsDir() {
			cmd.Dir = workingDir
		}
	}

	// Inherit PATH and standard environment
	cmd.Env = os.Environ()

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error al abrir stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error al abrir stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("fallo al iniciar comando compose: %w", err)
	}

	var wg sync.WaitGroup
	var capturedOutput strings.Builder
	var outputMu sync.Mutex

	scanStream := func(reader io.Reader) {
		defer wg.Done()
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			outputMu.Lock()
			capturedOutput.WriteString(line)
			capturedOutput.WriteString("\n")
			outputMu.Unlock()

			if onLog != nil {
				onLog(line)
			}
		}
	}

	wg.Add(2)
	go scanStream(stdoutPipe)
	go scanStream(stderrPipe)

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		outputMu.Lock()
		tail := capturedOutput.String()
		outputMu.Unlock()
		if len(tail) > 500 {
			tail = tail[len(tail)-500:]
		}
		if tail != "" {
			return fmt.Errorf("error en docker compose: %v (Detalle: %s)", err, strings.TrimSpace(tail))
		}
		return fmt.Errorf("error en docker compose: %w", err)
	}

	return nil
}

// ResolveComposeFile searches for standard compose file names in a directory or given configFile
func ResolveComposeFile(workingDir, configFile string) (string, bool) {
	if configFile != "" {
		// Handle comma-separated list of files (e.g. from com.docker.compose.project.config_files)
		paths := strings.Split(configFile, ",")
		for _, p := range paths {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
				return p, true
			}
			// If configFile is relative, try joining with workingDir
			if workingDir != "" && !filepath.IsAbs(p) {
				candidate := filepath.Join(workingDir, p)
				if fi, err := os.Stat(candidate); err == nil && !fi.IsDir() {
					return candidate, true
				}
			}
		}
	}

	if workingDir != "" {
		candidates := []string{
			"compose.yaml",
			"compose.yml",
			"docker-compose.yml",
			"docker-compose.yaml",
		}
		for _, name := range candidates {
			candidate := filepath.Join(workingDir, name)
			if fi, err := os.Stat(candidate); err == nil && !fi.IsDir() {
				return candidate, true
			}
		}
	}

	return "", false
}

// DiscoverDockerfiles searches for Dockerfiles associated with a stack in its workingDir
func DiscoverDockerfiles(workingDir string, composeContent string, serviceNames []string) []stackdomain.DockerfileInfo {
	if strings.TrimSpace(workingDir) == "" {
		return nil
	}

	cleanWorkingDir := filepath.Clean(workingDir)
	var results []stackdomain.DockerfileInfo
	seen := make(map[string]bool)

	addCandidate := func(candidatePath, serviceName string) {
		candidatePath = filepath.Clean(candidatePath)
		if seen[candidatePath] {
			return
		}
		fi, err := os.Stat(candidatePath)
		if err != nil || fi.IsDir() {
			return
		}

		data, err := os.ReadFile(candidatePath)
		if err != nil {
			return
		}

		seen[candidatePath] = true
		relName, err := filepath.Rel(cleanWorkingDir, candidatePath)
		if err != nil || strings.HasPrefix(relName, "..") {
			relName = filepath.Base(candidatePath)
		}

		results = append(results, stackdomain.DockerfileInfo{
			Name:         relName,
			Path:         candidatePath,
			ServiceName:  serviceName,
			Content:      string(data),
			ExistsOnDisk: true,
		})
	}

	// 1. Standard root candidate names
	rootCandidates := []string{
		"Dockerfile",
		"Containerfile",
		"Dockerfile.dev",
		"Dockerfile.prod",
		"Dockerfile.local",
	}
	for _, name := range rootCandidates {
		addCandidate(filepath.Join(cleanWorkingDir, name), "")
	}

	// 2. Service subdirectories
	for _, svc := range serviceNames {
		svc = strings.TrimSpace(svc)
		if svc == "" {
			continue
		}
		addCandidate(filepath.Join(cleanWorkingDir, svc, "Dockerfile"), svc)
		addCandidate(filepath.Join(cleanWorkingDir, svc, "Containerfile"), svc)
		addCandidate(filepath.Join(cleanWorkingDir, svc, "Dockerfile.dev"), svc)
		addCandidate(filepath.Join(cleanWorkingDir, "docker", svc, "Dockerfile"), svc)
		addCandidate(filepath.Join(cleanWorkingDir, "services", svc, "Dockerfile"), svc)
	}

	// 3. Shallow walk in workingDir (up to depth 2)
	ignoreDirs := map[string]bool{
		".git":         true,
		"node_modules": true,
		"vendor":       true,
		".cache":       true,
		"dist":         true,
		"build":        true,
		".idea":        true,
		".vscode":      true,
		"tmp":          true,
		"temp":         true,
	}

	_ = filepath.Walk(cleanWorkingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			base := info.Name()
			if ignoreDirs[base] {
				return filepath.SkipDir
			}
			rel, relErr := filepath.Rel(cleanWorkingDir, path)
			if relErr == nil && strings.Count(rel, string(os.PathSeparator)) >= 2 {
				return filepath.SkipDir
			}
			return nil
		}

		baseLower := strings.ToLower(info.Name())
		if strings.HasPrefix(baseLower, "dockerfile") ||
			strings.HasPrefix(baseLower, "containerfile") ||
			strings.HasSuffix(baseLower, ".dockerfile") {
			// Infer service if path matches a known service name
			inferredSvc := ""
			for _, svc := range serviceNames {
				if strings.Contains(path, svc) {
					inferredSvc = svc
					break
				}
			}
			addCandidate(path, inferredSvc)
		}
		return nil
	})

	return results
}

