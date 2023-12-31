package home

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication"

	homeHandler "github.com/peteraba/go-htmx-playground/pkg/home/handler"
)

func Setup[T authentication.Ctx](app *fiber.App, interceptor *authentication.Interceptor[T], logger *slog.Logger) {
	hHandler := homeHandler.NewHome(interceptor, logger)
	app.Use("/", adaptor.HTTPMiddleware(interceptor.CheckAuthentication()))
	app.Get("/", hHandler.Get)
}
