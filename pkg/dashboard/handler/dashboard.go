package handler

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/lib/auth"
	"github.com/peteraba/go-htmx-playground/pkg/dashboard/view"
)

type Dashboard struct {
	logger *slog.Logger
}

func NewDashboard(logger *slog.Logger) Dashboard {
	return Dashboard{
		logger: logger,
	}
}

func (d Dashboard) Get(c *fiber.Ctx) error {
	c.Accepts(fiber.MIMETextHTML)
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	userName := auth.GetString(c.Context(), auth.Name)
	component := view.Dashboard(userName)
	err := component.Render(c.Context(), c.Response().BodyWriter())

	return err
}
