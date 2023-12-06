package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/lib/htmx"
)

type Colors struct {
	themes []string
}

func NewColors() Colors {
	return Colors{
		themes: []string{
			"light",
			"dark",
			"cupcake",
			"bumblebee",
			"emerald",
			"corporate",
			"synthwave",
			"retro",
			"cyberpunk",
			"valentine",
			"halloween",
			"garden",
			"forest",
			"aqua",
			"lofi",
			"pastel",
			"fantasy",
			"wireframe",
			"black",
			"luxury",
			"dracula",
			"cmyk",
			"autumn",
			"business",
			"acid",
			"lemonade",
			"night",
			"coffee",
			"winter",
			"dim",
			"nord",
			"sunset",
		},
	}
}

func (h Colors) Get(c *fiber.Ctx) error {
	bind := fiber.Map{"Path": c.Path(), "Url": c.BaseURL(), "Themes": h.themes}

	if htmx.IsHx(c.GetReqHeaders()) {
		return c.Render("templates/colors", bind)
	}
	// Render index
	return c.Render("templates/colors", bind, "templates/layout")
}
