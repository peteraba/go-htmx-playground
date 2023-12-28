package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/pkg/home/view"
)

type Home struct{}

func NewHome() Home {
	return Home{}
}

func (h Home) Get(c *fiber.Ctx) error {
	component := view.Home()
	err := component.Render(c.Context(), c.Response().BodyWriter())

	return err
}
