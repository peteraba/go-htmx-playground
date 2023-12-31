package handler

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication"

	"github.com/peteraba/go-htmx-playground/pkg/home/view"
)

type Home[T authentication.Ctx] struct {
	interceptor *authentication.Interceptor[T]
	logger      *slog.Logger
}

func NewHome[T authentication.Ctx](interceptor *authentication.Interceptor[T], logger *slog.Logger) Home[T] {
	return Home[T]{
		interceptor: interceptor,
		logger:      logger,
	}
}

func (h Home[T]) Get(c *fiber.Ctx) error {
	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

	isAuthenticated := h.interceptor.Context(c.Context()).IsAuthenticated()

	component := view.Home(isAuthenticated)
	err := component.Render(c.Context(), c.Response().BodyWriter())

	return err
}
