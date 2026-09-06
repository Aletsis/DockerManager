package terminal

// StartResult holds the metadata returned when a PTY terminal is started
type StartResult struct {
	SessionID string `json:"sessionId"`
	Rows      uint   `json:"rows"`
	Cols      uint   `json:"cols"`
}
