package dashboard

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication"

	"github.com/peteraba/go-htmx-playground/pkg/dashboard/handler"
)

func Setup[T authentication.Ctx](app *fiber.App, interceptor *authentication.Interceptor[T], userInfoRetriever func(ctx context.Context) ([]byte, error), logger *slog.Logger) {
	dHandler := handler.NewDashboard[T](interceptor, userInfoRetriever, logger)
	app.Use("/dashboard", adaptor.HTTPMiddleware(interceptor.RequireAuthentication()))
	app.Get("/dashboard", dHandler.Get)
}
