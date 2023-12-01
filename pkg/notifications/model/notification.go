package model

type NotificationType string

const (
	INFO    NotificationType = "info"
	SUCCESS NotificationType = "success"
	WARNING NotificationType = "warning"
	ERROR   NotificationType = "error"
)

type Notification struct {
	Target  string           `json:"target"`
	Type    NotificationType `json:"type"`
	Message string           `json:"message"`
}
