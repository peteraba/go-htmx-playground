package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/lib/htmx"
)

type Home struct{}

func NewHome() Home {
	return Home{}
}

func (h Home) Get(c *fiber.Ctx) error {
	bind := fiber.Map{"Path": c.Path(), "Url": c.BaseURL()}

	if htmx.IsHx(c.GetReqHeaders()) {
		return c.Render("templates/home", bind)
	}
	// Render index
	return c.Render("templates/home", bind, "templates/layout")
}
