package service

import (
	"log/slog"

	"github.com/peteraba/go-htmx-playground/pkg/notifications/model"
)

type Notifier struct {
	logger          *slog.Logger
	sseChannelsByIP map[string][]chan model.Notification
}

func NewNotifier(logger *slog.Logger) *Notifier {
	return &Notifier{
		logger:          logger,
		sseChannelsByIP: make(map[string][]chan model.Notification),
	}
}

func (n *Notifier) Sub(ip string) chan model.Notification {
	sseChannel := make(chan model.Notification, 1)

	n.sseChannelsByIP[ip] = append(n.sseChannelsByIP[ip], sseChannel)

	return sseChannel
}

// nolint: varnamelen
func (n *Notifier) Unsub(ip string, sseChannel chan model.Notification) {
	for i, channel := range n.sseChannelsByIP[ip] {
		if channel != sseChannel {
			continue
		}

		n.sseChannelsByIP[ip] = append(n.sseChannelsByIP[ip][:i], n.sseChannelsByIP[ip][i+1:]...)

		close(channel)
	}
}

func (n *Notifier) broadcast(nType model.NotificationType, message string) {
	for _, sseChannels := range n.sseChannelsByIP {
		for _, sseChannel := range sseChannels {
			sseChannel <- model.Notification{
				Type:    nType,
				Message: message,
			}
		}
	}
}

func (n *Notifier) broadcastByIP(nType model.NotificationType, message, targetIP string) {
	sseChannels, ok := n.sseChannelsByIP[targetIP]
	if !ok {
		n.logger.Error("SSE channel not found. ip", slog.String("ip", targetIP))
	}

	for _, sseChannel := range sseChannels {
		sseChannel <- model.Notification{
			Type:    nType,
			Message: message,
		}
	}
}

func (n *Notifier) Info(message, targetIP string) {
	n.broadcastByIP(model.INFO, message, targetIP)
}

func (n *Notifier) Success(message, targetIP string) {
	n.broadcastByIP(model.SUCCESS, message, targetIP)
}

func (n *Notifier) Warning(message, targetIP string) {
	n.broadcastByIP(model.WARNING, message, targetIP)
}

func (n *Notifier) Error(message, targetIP string) {
	n.broadcastByIP(model.ERROR, message, targetIP)
}

func (n *Notifier) Reload() {
	n.broadcast(model.RELOAD, "")
}
