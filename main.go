package main

import (
	"embed"
	"log/slog"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"

	colorsHandler "github.com/peteraba/go-htmx-playground/pkg/colors/handler"
	filmsHandler "github.com/peteraba/go-htmx-playground/pkg/films/handler"
	"github.com/peteraba/go-htmx-playground/pkg/films/repository"
	homeHandler "github.com/peteraba/go-htmx-playground/pkg/home/handler"
	notificationsHandler "github.com/peteraba/go-htmx-playground/pkg/notifications/handler"
	notificationsService "github.com/peteraba/go-htmx-playground/pkg/notifications/service"
	"github.com/peteraba/go-htmx-playground/pkg/view"
)

//go:embed assets/*
var assetsFS embed.FS

func main() {
	const maxListLength = 10

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	notifier := notificationsService.NewNotifier(logger)
	app := fiber.New(fiber.Config{
		Views:     view.NewEngine(),
		Immutable: true,
	})

	setupMiddleware(app, logger)

	addSseHandler(app, logger, notifier)
	addHomeHandler(app)
	addColorHandlers(app)
	addFilmHandlers(app, logger, notifier, maxListLength)
	addStaticHandler(app)

	err := app.Listen(":8000")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func setupMiddleware(app *fiber.App, logger *slog.Logger) *slog.Logger {
	app.Use(recover.New())

	app.Get("/metrics", monitor.New(monitor.Config{Title: "go|htmx Metrics Page"}))
	app.Use(slogfiber.New(logger))
	app.Use(idempotency.New())

	return logger
}

func addSseHandler(app *fiber.App, logger *slog.Logger, notifier *notificationsService.Notifier) *notificationsService.Notifier {
	handler := notificationsHandler.NewSSE(logger, notifier)
	app.Get("/messages", handler.ServeMessages)

	return notifier
}

func addHomeHandler(app *fiber.App) {
	hHandler := homeHandler.NewHome()
	app.Get("/", hHandler.Get)
}

func addColorHandlers(app *fiber.App) {
	cHandler := colorsHandler.NewColors()
	app.Get("/colors", cHandler.Get)
}

func addFilmHandlers(app *fiber.App, logger *slog.Logger, notifier *notificationsService.Notifier, maxListLength int) {
	repo := repository.NewFilmRepo(logger, maxListLength)

	fHandler := filmsHandler.NewFilm(logger, repo, notifier, maxListLength)
	app.Get("/films", fHandler.List)
	app.Post("/films", fHandler.Create)
	app.Post("/films-delete", fHandler.DeleteForm)
	app.Delete("/films", fHandler.Delete)
	app.Post("/generators/films/:num<min(5);max(50)>", fHandler.Generate)

	dHandler := filmsHandler.NewDirector(repo, maxListLength)
	app.Get("/directors", dHandler.List)
}

func addStaticHandler(app *fiber.App) {
	app.Use(filesystem.New(filesystem.Config{
		Root: http.FS(assetsFS),
		// PathPrefix:   "/assets",
		Browse:       false,
		NotFoundFile: "404.html",
		MaxAge:       3600,
	}))
}
