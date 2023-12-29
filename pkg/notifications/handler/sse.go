package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/lib/log"
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

func (s SSE) handleEvent(w *bufio.Writer, sseEvent model.Notification) bool {
	var (
		bolB []byte
		size int
		err  error
	)

	s.logger.Info("1")

	bolB, err = json.Marshal(sseEvent)
	if err != nil {
		s.logger.Error("Error while marshaling JSON.", log.Err(err))

		return false
	}

	s.logger.Info("2")

	size, err = fmt.Fprintf(w, "data: %s\n\n", string(bolB))
	if err != nil {
		s.logger.Error("Error while writing buffer.", log.Err(err))

		return false
	}

	s.logger.Info("3")
	s.logger.With("type", sseEvent.Type, "bytes", size).Info("Message sent.")

	err = w.Flush()
	if err != nil {
		// Refreshing page in web browser will establish a new
		// SSE connection, but only (the last) one is alive, so
		// dead connections must be closed here.
		s.logger.Error("Error while flushing. Closing HTTP connection.", log.Err(err))

		return false
	}

	s.logger.Info("4")
	s.logger.Info("Messages flushed.")

	return true
}

func (s SSE) ServeMessages(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	sseChannel := s.notifier.Sub(c.IP())

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer func() {
			s.logger.With("ip", c.IP()).Info("SSE connection closed.")
			s.notifier.Unsub(c.IP(), sseChannel)
		}()

		s.logger.Info("Broadcasting messages is ready.")
		s.logger.Info("A")

		sseChannel <- model.Notification{Type: model.RELOAD, Message: ""}

		s.logger.Info("B")

		for {
			if !s.handleEvent(w, <-sseChannel) {
				return
			}
		}
	})

	return nil
}
