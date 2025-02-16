package models

type ResponseErr struct {
	Error      string `json:"error,omitempty"`
	StatusCode int    `json:"status_code,omitempty"`
}
