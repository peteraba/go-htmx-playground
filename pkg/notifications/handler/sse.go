package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/pkg/notifications/model"
	"github.com/peteraba/go-htmx-playground/pkg/notifications/service"
)

type SSE struct {
	logger   *slog.Logger
	notifier *service.Notifier
}

func NewSSE(logger *slog.Logger, notifier *service.Notifier) SSE {
	return SSE{
		logger:   logger,
		notifier: notifier,
	}
}

func (s SSE) ServeMessages(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	sseChannel := s.notifier.Sub(c.IP())

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		s.logger.Info("Broadcasting messages is ready.")
		var (
			sseEvent model.Notification
			bolB     []byte
			size     int
			err      error
		)

		sseChannel <- model.Notification{Type: model.RELOAD, Message: ""}

		for {
			sseEvent = <-sseChannel

			bolB, err = json.Marshal(sseEvent)
			if err != nil {
				s.logger.With("err", err).Error("Error while marshaling JSON.")

				return
			}

			size, err = fmt.Fprintf(w, "data: %s\n\n", string(bolB))
			if err != nil {
				s.logger.With("err", err).Error("Error while writing buffer.")

				return
			}
			s.logger.With("type", sseEvent.Type, "bytes", size).Info("Message sent.")

			err = w.Flush()
			if err != nil {
				// Refreshing page in web browser will establish a new
				// SSE connection, but only (the last) one is alive, so
				// dead connections must be closed here.
				s.logger.With("err", err).Error("Error while flushing. Closing http connection.")

				return
			}
		}
	})

	return nil
}
