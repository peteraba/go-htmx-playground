package handler

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/pkg/home/view"
)

type Home struct {
	logger *slog.Logger
}

func NewHome(logger *slog.Logger) Home {
	return Home{
		logger: logger,
	}
}

func (h Home) Get(c *fiber.Ctx) error {
	c.Accepts(fiber.MIMETextHTML)
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	component := view.Home()
	err := component.Render(c.Context(), c.Response().BodyWriter())

	return err
}
