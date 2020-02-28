package handlers

// User for S3 access
type User struct {
	Username string `json:"user"`
	Basepath string `json:"path"`
}
