package system

// Overview represents host and Docker engine metrics
type Overview struct {
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
