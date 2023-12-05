package model

type NotificationType string

const (
	INFO    NotificationType = "info"
	SUCCESS NotificationType = "success"
	WARNING NotificationType = "warning"
	ERROR   NotificationType = "error"
	RELOAD  NotificationType = "reload"
)

type Notification struct {
	Type    NotificationType `json:"type"`
	Message string           `json:"message"`
}
