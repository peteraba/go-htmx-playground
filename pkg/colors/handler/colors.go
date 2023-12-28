package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/pkg/colors/view"
)

type Colors struct {
	themes view.Themes
}

func NewColors() Colors {
	return Colors{
		themes: view.NewThemes(),
	}
}

func (h Colors) Get(c *fiber.Ctx) error {
	component := h.themes.Colors()
	err := component.Render(c.Context(), c.Response().BodyWriter())

	return err
}
