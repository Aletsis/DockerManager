package container

// State represents the lifecycle state of a container
type State string

const (
	StateRunning    State = "running"
	StatePaused     State = "paused"
	StateExited     State = "exited"
	StateRestarting State = "restarting"
	StateDead       State = "dead"
	StateCreated    State = "created"
)

// PortMapping represents network port exposure between host and container
type PortMapping struct {
	IP          string `json:"ip"`
	PrivatePort uint16 `json:"privatePort"`
	PublicPort  uint16 `json:"publicPort"`
	Type        string `json:"type"`
}

// NetworkInfo represents IP and network configuration for a container
type NetworkInfo struct {
	NetworkName string   `json:"networkName"`
	NetworkID   string   `json:"networkId"`
	IPAddress   string   `json:"ipAddress"`
	Gateway     string   `json:"gateway"`
	MacAddress  string   `json:"macAddress"`
	Aliases     []string `json:"aliases,omitempty"`
}

// Stats represents real-time resource utilization
type Stats struct {
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
	Pids             uint64  `json:"pids"`
}

// CreateSpec contains specifications required to create a container
type CreateSpec struct {
	Image         string   `json:"image"`
	Name          string   `json:"name,omitempty"`
	Cmd           []string `json:"cmd,omitempty"`
	Ports         []string `json:"ports,omitempty"`
	Volumes       []string `json:"volumes,omitempty"`
	Env           []string `json:"env,omitempty"`
	RestartPolicy string   `json:"restartPolicy,omitempty"`
	AutoStart     bool     `json:"autoStart,omitempty"`
}

// CreateResult contains the outcome of a container creation
type CreateResult struct {
	ID       string   `json:"id"`
	Warnings []string `json:"warnings,omitempty"`
}
