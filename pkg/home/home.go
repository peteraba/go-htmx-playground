package home

import (
	"github.com/gofiber/fiber/v2"

	homeHandler "github.com/peteraba/go-htmx-playground/pkg/home/handler"
)

func Setup(app *fiber.App) {
	hHandler := homeHandler.NewHome()
	app.Get("/", hHandler.Get)
}
