package docker

// PortMapping represents a published or exposed container port
type PortMapping struct {
	IP          string `json:"ip"`
	PrivatePort uint16 `json:"privatePort"`
	PublicPort  uint16 `json:"publicPort"`
	Type        string `json:"type"`
}

// ContainerInfo represents detailed information about a Docker container
type ContainerInfo struct {
	ID         string        `json:"id"`
	ShortID    string        `json:"shortId"`
	Names      []string      `json:"names"`
	Name       string        `json:"name"`
	Image      string        `json:"image"`
	ImageID    string        `json:"imageId"`
	Command    string        `json:"command"`
	Created    int64         `json:"created"`
	State      string        `json:"state"`
	Status     string        `json:"status"`
	Ports            []PortMapping     `json:"ports"`
	SizeRw           int64             `json:"sizeRw"`
	SizeRootFs       int64             `json:"sizeRootFs"`
	Labels           map[string]string `json:"labels,omitempty"`
	ComposeProject   string            `json:"composeProject,omitempty"`
	ComposeService   string            `json:"composeService,omitempty"`
	ComposeWorkingDir string           `json:"composeWorkingDir,omitempty"`
	ComposeConfigFile string           `json:"composeConfigFile,omitempty"`
}

// SystemOverview provides high-level metrics about the Docker daemon
type SystemOverview struct {
	Containers        int    `json:"containers"`
	ContainersRunning int    `json:"containersRunning"`
	ContainersPaused  int    `json:"containersPaused"`
	ContainersStopped int    `json:"containersStopped"`
	Images            int    `json:"images"`
	ServerVersion     string `json:"serverVersion"`
	OperatingSystem   string `json:"operatingSystem"`
	NCPU              int    `json:"ncpu"`
	MemTotal          int64  `json:"memTotal"`
}

// ContainerStats represents the current resource utilization of a container
type ContainerStats struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	CPUPercentage    float64 `json:"cpuPercentage"`
	MemoryUsage      uint64  `json:"memoryUsage"`
	MemoryLimit      uint64  `json:"memoryLimit"`
	MemoryPercentage float64 `json:"memoryPercentage"`
	NetworkRx        uint64  `json:"networkRx"`
	NetworkTx        uint64  `json:"networkTx"`
	BlockRead        uint64  `json:"blockRead"`
	BlockWrite       uint64  `json:"blockWrite"`
	PIDs             uint64  `json:"pids"`
}

// ImageInfo represents detailed information about a Docker image
type ImageInfo struct {
	ID          string   `json:"id"`
	ShortID     string   `json:"shortId"`
	Repository  string   `json:"repository"`
	Tag         string   `json:"tag"`
	RepoTags    []string `json:"repoTags"`
	Created     int64    `json:"created"`
	Size        int64    `json:"size"`
	SharedSize  int64    `json:"sharedSize"`
	Containers  int64    `json:"containers"`
	InUse       bool     `json:"inUse"`
	IsDangling  bool     `json:"isDangling"`
}

// DiskUsageSummary summarizes disk usage metrics for images and dangling assets
type DiskUsageSummary struct {
	TotalImages     int   `json:"totalImages"`
	TotalSize       int64 `json:"totalSize"`
	DanglingCount   int   `json:"danglingCount"`
	DanglingSize    int64 `json:"danglingSize"`
	ReclaimableSize int64 `json:"reclaimableSize"`
}

// PruneResult represents the outcome of an image prune cleanup operation
type PruneResult struct {
	ImagesDeleted  []string `json:"imagesDeleted"`
	SpaceReclaimed uint64   `json:"spaceReclaimed"`
}

// PullProgressEvent reports streaming progress for image pull operations
type PullProgressEvent struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Progress string `json:"progress"`
	Current  int64  `json:"current"`
	Total    int64  `json:"total"`
	Error    string `json:"error,omitempty"`
}

// CreateContainerRequest defines the parameters to create and start a new container
type CreateContainerRequest struct {
	Image         string   `json:"image"`
	Name          string   `json:"name"`
	Ports         []string `json:"ports"`
	Volumes       []string `json:"volumes"`
	Env           []string `json:"env"`
	RestartPolicy string   `json:"restartPolicy"`
	AutoStart     bool     `json:"autoStart"`
}

// CreateContainerResult contains details of a newly created container
type CreateContainerResult struct {
	ID       string   `json:"id"`
	Warnings []string `json:"warnings"`
}

