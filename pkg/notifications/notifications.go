package notifications

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	notificationsHandler "github.com/peteraba/go-htmx-playground/pkg/notifications/handler"
	notificationsService "github.com/peteraba/go-htmx-playground/pkg/notifications/service"
)

func Setup(app *fiber.App, logger *slog.Logger, notifier *notificationsService.Notifier) {
	handler := notificationsHandler.NewSSE(logger, notifier)
	app.Get("/messages", handler.ServeMessages)
}
