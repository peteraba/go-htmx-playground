package dashboard

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/pkg/dashboard/handler"
)

func Setup(app *fiber.App, requireAuthHandler fiber.Handler, logger *slog.Logger) {
	dHandler := handler.NewDashboard(logger)
	app.Use("/dashboard", requireAuthHandler)
	app.Get("/dashboard", dHandler.Get)
}
