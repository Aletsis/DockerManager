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
	Ports      []PortMapping `json:"ports"`
	SizeRw     int64         `json:"sizeRw"`
	SizeRootFs int64         `json:"sizeRootFs"`
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
