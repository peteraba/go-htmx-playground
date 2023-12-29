package main

import (
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/pkg/colors"
	"github.com/peteraba/go-htmx-playground/pkg/films"
	"github.com/peteraba/go-htmx-playground/pkg/home"
	"github.com/peteraba/go-htmx-playground/pkg/notifications"
	notificationsService "github.com/peteraba/go-htmx-playground/pkg/notifications/service"
	"github.com/peteraba/go-htmx-playground/pkg/server"
)

var Version = "development" // nolint:gochecknoglobals

func main() {
	const maxListLength = 10

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	notifier := notificationsService.NewNotifier(logger)
	//nolint: exhaustruct
	app := fiber.New(fiber.Config{
		Immutable: true,
	})

	server.Setup(app, logger, Version)
	notifications.Setup(app, logger, notifier)
	home.Setup(app)
	colors.Setup(app)
	films.Setup(app, logger, notifier, maxListLength, Version)
	server.AddStaticHandler(app)

	err := app.Listen(":8000")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
