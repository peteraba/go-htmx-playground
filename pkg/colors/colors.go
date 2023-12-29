package colors

import (
	"github.com/gofiber/fiber/v2"

	colorsHandler "github.com/peteraba/go-htmx-playground/pkg/colors/handler"
)

func Setup(app *fiber.App) {
	cHandler := colorsHandler.NewColors()
	app.Get("/colors", cHandler.Get)
}
