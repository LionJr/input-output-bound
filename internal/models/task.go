package models

type Status string

const (
	Pending Status = "pending"
	Running Status = "running"
	Done    Status = "done"
	Error   Status = "error"
)

type Task struct {
	Id     string `json:"id"`
	Status Status `json:"status"`
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}
