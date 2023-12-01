package service

import (
	"log"

	"github.com/peteraba/go-htmx-playground/pkg/notifications/model"
)

type Notifier struct {
	sseChannelsByIP map[string][]chan model.Notification
}

func NewNotifier() *Notifier {
	sseChannelsByIP := make(map[string][]chan model.Notification)
	return &Notifier{
		sseChannelsByIP: sseChannelsByIP,
	}
}

func (n *Notifier) Sub(ip string) chan model.Notification {
	sseChannel := make(chan model.Notification)

	n.sseChannelsByIP[ip] = append(n.sseChannelsByIP[ip], sseChannel)

	return sseChannel
}

func (n *Notifier) broadcastByIP(nType model.NotificationType, message, targetIP string) {
	sseChannels, ok := n.sseChannelsByIP[targetIP]
	if !ok {
		log.Println("SSE channel not found")
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
