package docker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	stackdomain "dockermanager/internal/domain/stack"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

// StackRepository implements stackdomain.Repository using Docker SDK and ComposeRunner
type StackRepository struct {
	client *Client
	runner *ComposeRunner
}

// NewStackRepository creates a new Docker StackRepository
func NewStackRepository(client *Client) *StackRepository {
	return &StackRepository{
		client: client,
		runner: NewComposeRunner(),
	}
}

// StartStack starts all stopped containers in a compose project
func (r *StackRepository) StartStack(ctx context.Context, projectName string) error {
	projectName = strings.TrimSpace(projectName)
	if projectName == "" {
		return fmt.Errorf("el nombre del proyecto compose no puede estar vacío")
	}

	filterArgs := filters.NewArgs()
	filterArgs.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))

	containers, err := r.client.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filterArgs})
	if err != nil {
		return fmt.Errorf("error al listar contenedores del stack %s: %w", projectName, err)
	}
	if len(containers) == 0 {
		return fmt.Errorf("no se encontraron contenedores para el stack %q", projectName)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []string

	for _, c := range containers {
		if c.State == "running" {
			continue
		}
		wg.Add(1)
		go func(cid string) {
			defer wg.Done()
			if err := r.client.cli.ContainerStart(ctx, cid, container.StartOptions{}); err != nil {
				mu.Lock()
				short := cid
				if len(short) > 12 {
					short = short[:12]
				}
				errs = append(errs, fmt.Sprintf("%s: %v", short, err))
				mu.Unlock()
			}
		}(c.ID)
	}
	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("errores al iniciar el stack: %s", strings.Join(errs, "; "))
	}
	return nil
}

// StopStack stops all running containers in a compose project
func (r *StackRepository) StopStack(ctx context.Context, projectName string) error {
	projectName = strings.TrimSpace(projectName)
	if projectName == "" {
		return fmt.Errorf("el nombre del proyecto compose no puede estar vacío")
	}

	filterArgs := filters.NewArgs()
	filterArgs.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))

	containers, err := r.client.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filterArgs})
	if err != nil {
		return fmt.Errorf("error al listar contenedores del stack %s: %w", projectName, err)
	}
	if len(containers) == 0 {
		return fmt.Errorf("no se encontraron contenedores para el stack %q", projectName)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []string
	timeout := 15

	for _, c := range containers {
		if c.State != "running" {
			continue
		}
		wg.Add(1)
		go func(cid string) {
			defer wg.Done()
			if err := r.client.cli.ContainerStop(ctx, cid, container.StopOptions{Timeout: &timeout}); err != nil {
				mu.Lock()
				short := cid
				if len(short) > 12 {
					short = short[:12]
				}
				errs = append(errs, fmt.Sprintf("%s: %v", short, err))
				mu.Unlock()
			}
		}(c.ID)
	}
	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("errores al detener el stack: %s", strings.Join(errs, "; "))
	}
	return nil
}

// RestartStack restarts all containers in a compose project
func (r *StackRepository) RestartStack(ctx context.Context, projectName string) error {
	projectName = strings.TrimSpace(projectName)
	if projectName == "" {
		return fmt.Errorf("el nombre del proyecto compose no puede estar vacío")
	}

	filterArgs := filters.NewArgs()
	filterArgs.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))

	containers, err := r.client.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filterArgs})
	if err != nil {
		return fmt.Errorf("error al listar contenedores del stack %s: %w", projectName, err)
	}
	if len(containers) == 0 {
		return fmt.Errorf("no se encontraron contenedores para el stack %q", projectName)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []string
	timeout := 15

	for _, c := range containers {
		wg.Add(1)
		go func(cid string) {
			defer wg.Done()
			if err := r.client.cli.ContainerRestart(ctx, cid, container.StopOptions{Timeout: &timeout}); err != nil {
				mu.Lock()
				short := cid
				if len(short) > 12 {
					short = short[:12]
				}
				errs = append(errs, fmt.Sprintf("%s: %v", short, err))
				mu.Unlock()
			}
		}(c.ID)
	}
	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("errores al reiniciar el stack: %s", strings.Join(errs, "; "))
	}
	return nil
}

// GetDefaultStackDirectory returns ~/.dockermanager/stacks/<projectName>
func (r *StackRepository) GetDefaultStackDirectory(projectName string) (string, error) {
	projectName = strings.TrimSpace(projectName)
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error al obtener directorio home del usuario: %w", err)
	}
	dir := filepath.Join(home, ".dockermanager", "stacks")
	if projectName != "" {
		dir = filepath.Join(dir, projectName)
	}
	return dir, nil
}

// SaveComposeFile writes or updates the compose file content on disk
func (r *StackRepository) SaveComposeFile(ctx context.Context, workingDir, configFile, content string) (string, error) {
	workingDir = strings.TrimSpace(workingDir)
	configFile = strings.TrimSpace(configFile)

	targetPath := configFile
	if targetPath == "" {
		if workingDir == "" {
			return "", fmt.Errorf("se debe especificar workingDir o configFile para guardar el archivo")
		}
		targetPath = filepath.Join(workingDir, "compose.yaml")
	} else if !filepath.IsAbs(targetPath) && workingDir != "" {
		targetPath = filepath.Join(workingDir, targetPath)
	}

	dir := filepath.Dir(targetPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("error al crear directorio %s: %w", dir, err)
	}

	if err := os.WriteFile(targetPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("error al escribir archivo %s: %w", targetPath, err)
	}

	return targetPath, nil
}

// GetComposeFile retrieves the compose file content, metadata, and discovered Dockerfiles for a stack
func (r *StackRepository) GetComposeFile(ctx context.Context, projectName, workingDir, configFile string) (*stackdomain.ComposeFileInfo, error) {
	projectName = strings.TrimSpace(projectName)
	workingDir = strings.TrimSpace(workingDir)
	configFile = strings.TrimSpace(configFile)

	var serviceNames []string

	// Inspect containers to discover working_dir, config_files, and service names from labels
	if projectName != "" {
		filterArgs := filters.NewArgs()
		filterArgs.Add("label", fmt.Sprintf("com.docker.compose.project=%s", projectName))
		containers, err := r.client.cli.ContainerList(ctx, container.ListOptions{All: true, Filters: filterArgs})
		if err == nil && len(containers) > 0 {
			svcMap := make(map[string]bool)
			for _, c := range containers {
				if c.Labels != nil {
					if workingDir == "" && c.Labels["com.docker.compose.project.working_dir"] != "" {
						workingDir = c.Labels["com.docker.compose.project.working_dir"]
					}
					if configFile == "" && c.Labels["com.docker.compose.project.config_files"] != "" {
						configFile = c.Labels["com.docker.compose.project.config_files"]
					}
					if svc := c.Labels["com.docker.compose.service"]; svc != "" && !svcMap[svc] {
						svcMap[svc] = true
						serviceNames = append(serviceNames, svc)
					}
				}
			}
		}
	}

	info := &stackdomain.ComposeFileInfo{
		ProjectName:  projectName,
		WorkingDir:   workingDir,
		ConfigFile:   configFile,
		Content:      "",
		ExistsOnDisk: false,
		Dockerfiles:  []stackdomain.DockerfileInfo{},
	}

	resolvedPath, exists := ResolveComposeFile(workingDir, configFile)
	if exists {
		data, err := os.ReadFile(resolvedPath)
		if err == nil {
			info.ConfigFile = resolvedPath
			info.Content = string(data)
			info.ExistsOnDisk = true
			if info.WorkingDir == "" {
				info.WorkingDir = filepath.Dir(resolvedPath)
			}
		}
	}

	// If not found on disk, but workingDir was provided, check if default compose.yaml can be suggested
	if workingDir != "" && info.ConfigFile == "" {
		info.ConfigFile = filepath.Join(workingDir, "compose.yaml")
	}

	effectiveWorkingDir := info.WorkingDir
	if effectiveWorkingDir == "" && info.ConfigFile != "" {
		effectiveWorkingDir = filepath.Dir(info.ConfigFile)
	}

	// Automatically discover and load all Dockerfiles in the stack working directory
	if effectiveWorkingDir != "" {
		info.Dockerfiles = DiscoverDockerfiles(effectiveWorkingDir, info.Content, serviceNames)
	}

	return info, nil
}

// UpStack deploys a compose stack with 'docker compose up -d'
func (r *StackRepository) UpStack(ctx context.Context, req stackdomain.ComposeDeployRequest, onLog func(string)) error {
	projectName := strings.TrimSpace(req.ProjectName)
	if projectName == "" {
		return fmt.Errorf("el nombre del stack no puede estar vacío")
	}

	workingDir := strings.TrimSpace(req.WorkingDir)
	configFile := strings.TrimSpace(req.ConfigFile)

	// If content is provided, save it first
	if strings.TrimSpace(req.Content) != "" {
		savedPath, err := r.SaveComposeFile(ctx, workingDir, configFile, req.Content)
		if err != nil {
			return fmt.Errorf("error al guardar archivo compose antes de desplegar: %w", err)
		}
		configFile = savedPath
		if workingDir == "" {
			workingDir = filepath.Dir(savedPath)
		}
	} else if configFile == "" && workingDir != "" {
		// Attempt to resolve file
		if resolved, ok := ResolveComposeFile(workingDir, ""); ok {
			configFile = resolved
		}
	}

	args := []string{}
	if projectName != "" {
		args = append(args, "-p", projectName)
	}
	if configFile != "" {
		args = append(args, "-f", configFile)
	}
	args = append(args, "up", "-d")
	if req.RemoveOrphans {
		args = append(args, "--remove-orphans")
	}

	if onLog != nil {
		onLog(fmt.Sprintf("Iniciando despliegue de stack '%s'...", projectName))
		if configFile != "" {
			onLog(fmt.Sprintf("Archivo: %s", configFile))
		}
		if workingDir != "" {
			onLog(fmt.Sprintf("Directorio: %s", workingDir))
		}
	}

	err := r.runner.Run(ctx, workingDir, args, onLog)
	if err != nil {
		if onLog != nil {
			onLog(fmt.Sprintf("Fallo en despliegue: %v", err))
		}
		return err
	}

	if onLog != nil {
		onLog(fmt.Sprintf("Stack '%s' desplegado exitosamente.", projectName))
	}
	return nil
}

// DownStack tears down a compose stack with 'docker compose down'
func (r *StackRepository) DownStack(ctx context.Context, req stackdomain.ComposeDownRequest, onLog func(string)) error {
	projectName := strings.TrimSpace(req.ProjectName)
	if projectName == "" {
		return fmt.Errorf("el nombre del stack no puede estar vacío")
	}

	workingDir := strings.TrimSpace(req.WorkingDir)
	configFile := strings.TrimSpace(req.ConfigFile)

	// If configFile or workingDir not passed, attempt discovery from running containers
	if configFile == "" && workingDir == "" {
		info, _ := r.GetComposeFile(ctx, projectName, "", "")
		if info != nil {
			configFile = info.ConfigFile
			workingDir = info.WorkingDir
		}
	}

	args := []string{}
	if projectName != "" {
		args = append(args, "-p", projectName)
	}
	if configFile != "" {
		args = append(args, "-f", configFile)
	}
	args = append(args, "down")
	if req.RemoveVolumes {
		args = append(args, "--volumes")
	}

	if onLog != nil {
		onLog(fmt.Sprintf("Deteniendo y desmontando stack '%s'...", projectName))
	}

	err := r.runner.Run(ctx, workingDir, args, onLog)
	if err != nil {
		if onLog != nil {
			onLog(fmt.Sprintf("Fallo al desmontar stack: %v", err))
		}
		return err
	}

	if onLog != nil {
		onLog(fmt.Sprintf("Stack '%s' desmontado y eliminado exitosamente.", projectName))
	}
	return nil
}

