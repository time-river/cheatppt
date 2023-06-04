package api

type status string

const (
	FAILURE      = "failure"
	SUCCESS      = "success"
	UNAUTHORIZED = "unauthorized"
)

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Status  status      `json:"status"`
}
