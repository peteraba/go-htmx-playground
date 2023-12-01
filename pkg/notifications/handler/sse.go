package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/pkg/notifications/service"
)

type SSE struct {
	notifier *service.Notifier
}

func NewSSE(notifier *service.Notifier) SSE {
	return SSE{
		notifier: notifier,
	}
}

func (s SSE) Serve(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	sseChannel := s.notifier.Sub(c.IP())

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		log.Printf("SSE ready.")
		var (
			sseEvent interface{}
			bolB     []byte
			size     int
			err      error
		)
		for {
			sseEvent = <-sseChannel

			bolB, err = json.Marshal(sseEvent)
			if err != nil {
				log.Printf("Error while marshaling: %v. Closing http connection.", err)
				return
			}

			size, err = fmt.Fprintf(w, "data: %s\n\n", string(bolB))
			log.Printf("Message sent. (%d bytes)", size)

			err = w.Flush()
			if err != nil {
				// Refreshing page in web browser will establish a new
				// SSE connection, but only (the last) one is alive, so
				// dead connections must be closed here.
				log.Printf("Error while flushing: %v. Closing http connection.", err)

				return
			}
		}
	})

	return nil
}
