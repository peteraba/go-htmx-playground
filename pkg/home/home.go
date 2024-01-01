package home

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	homeHandler "github.com/peteraba/go-htmx-playground/pkg/home/handler"
)

func Setup(app *fiber.App, logger *slog.Logger) {
	hHandler := homeHandler.NewHome(logger)
	app.Get("/", hHandler.Get)
}
