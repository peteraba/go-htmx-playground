package main

import (
	"embed"
	"log/slog"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"

	"github.com/peteraba/go-htmx-playground/pkg/app/middleware"
	colorsHandler "github.com/peteraba/go-htmx-playground/pkg/colors/handler"
	filmsHandler "github.com/peteraba/go-htmx-playground/pkg/films/handler"
	"github.com/peteraba/go-htmx-playground/pkg/films/repository"
	filmsService "github.com/peteraba/go-htmx-playground/pkg/films/service"
	homeHandler "github.com/peteraba/go-htmx-playground/pkg/home/handler"
	notificationsHandler "github.com/peteraba/go-htmx-playground/pkg/notifications/handler"
	notificationsService "github.com/peteraba/go-htmx-playground/pkg/notifications/service"
)

var Version = "development" // nolint:gochecknoglobals

//go:embed assets/*
var assetsFS embed.FS

func main() {
	const maxListLength = 10

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	notifier := notificationsService.NewNotifier(logger)
	//nolint: exhaustruct
	app := fiber.New(fiber.Config{
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
	app.Use(middleware.Htmx(Version))

	//nolint: exhaustruct
	app.Get("/metrics", monitor.New(monitor.Config{Title: "go|htmx Metrics Page"}))
	app.Use(slogfiber.New(logger))
	// app.Use(idempotency.New())

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

	fService := filmsService.NewFilm(repo, logger)
	fHandler := filmsHandler.NewFilm(fService, maxListLength, notifier, logger, Version)
	app.Get("/films", fHandler.List)
	app.Post("/films", fHandler.Create)
	app.Post("/films-delete", fHandler.DeleteForm)
	app.Delete("/films", fHandler.Delete)
	app.Post("/generators/films/:num<min(5);max(50)>", fHandler.Generate)

	dHandler := filmsHandler.NewDirector(repo, maxListLength)
	app.Get("/directors", dHandler.List)
}

func addStaticHandler(app *fiber.App) {
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
