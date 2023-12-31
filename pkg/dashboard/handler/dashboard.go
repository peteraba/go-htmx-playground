package handler

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication"

	"github.com/peteraba/go-htmx-playground/pkg/dashboard/view"
)

type Dashboard[T authentication.Ctx] struct {
	interceptor       *authentication.Interceptor[T]
	userInfoRetriever func(ctx context.Context) ([]byte, error)
	logger            *slog.Logger
}

func NewDashboard[T authentication.Ctx](interceptor *authentication.Interceptor[T], userInfoRetriever func(ctx context.Context) ([]byte, error), logger *slog.Logger) Dashboard[T] {
	return Dashboard[T]{
		interceptor:       interceptor,
		userInfoRetriever: userInfoRetriever,
		logger:            logger,
	}
}

func (d Dashboard[T]) Get(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	data, err := d.userInfoRetriever(c.Context())
	if err != nil {
		d.logger.Error(err.Error())
		return c.Redirect("", fiber.StatusTemporaryRedirect)
	}
	// var data []byte
	// var err error

	component := view.Dashboard(data)
	err = component.Render(c.Context(), c.Response().BodyWriter())

	return err
}
