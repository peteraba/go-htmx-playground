package server

import (
	"embed"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"

	"github.com/peteraba/go-htmx-playground/pkg/server/middleware"
)

//go:embed assets/*
var assetsFS embed.FS

func Setup(app *fiber.App, logger *slog.Logger, version string) {
	app.Use(recover.New())
	app.Use(middleware.Htmx(version))

	//nolint: exhaustruct
	app.Get("/metrics", monitor.New(monitor.Config{Title: "go|htmx Metrics Page"}))
	app.Use(slogfiber.New(logger))
	app.Use(idempotency.New())
}

func AddStaticHandler(app *fiber.App) {
	const maxAge = 60 * 60

	//nolint: exhaustruct
	app.Use(filesystem.New(filesystem.Config{
		Root: http.FS(assetsFS),
		// PathPrefix:   "/assets",
		Browse:       false,
		NotFoundFile: "404.html",
		MaxAge:       maxAge,
	}))
}
