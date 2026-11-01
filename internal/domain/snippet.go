package domain

// Snippet is a saved shell command the user runs on one or many servers.
type Snippet struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	Description string `json:"description,omitempty"`
}
