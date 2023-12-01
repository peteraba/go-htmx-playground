package main

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	html "github.com/gofiber/template/html/v2"

	colorsHandler "github.com/peteraba/go-htmx-playground/pkg/colors/handler"
	filmsHandler "github.com/peteraba/go-htmx-playground/pkg/films/handler"
	"github.com/peteraba/go-htmx-playground/pkg/films/repository"
	homeHandler "github.com/peteraba/go-htmx-playground/pkg/home/handler"
	notificationsHandler "github.com/peteraba/go-htmx-playground/pkg/notifications/handler"
	notificationsService "github.com/peteraba/go-htmx-playground/pkg/notifications/service"
)

//go:embed views/*.html
var viewsFS embed.FS

//go:embed assets/*
var assetsFS embed.FS

func main() {
	const maxListLength = 5

	engine := html.NewFileSystem(http.FS(viewsFS), ".html")
	engine.
		AddFunc("str_slice", func(s []string) template.HTML {
			return template.HTML(strings.Join(s, ","))
		}).
		AddFunc("url", func(s string) template.URL {
			return template.URL(s)
		})

	app := fiber.New(fiber.Config{
		Views:     engine,
		Immutable: true,
	})

	//// See here: https://github.com/samber/slog-fiber
	//logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	//app.Use(slogfiber.New(logger))

	app.Use(recover.New())

	app.Use(idempotency.New())

	app.Get("/metrics", monitor.New(monitor.Config{Title: "go|htmx Metrics Page"}))

	repo := repository.NewFilmRepo(maxListLength)

	notifier := notificationsService.NewNotifier()
	sse := notificationsHandler.NewSSE(notifier)
	app.Get("/sse", sse.Serve)

	cHandler := colorsHandler.NewColors()
	app.Get("/colors", cHandler.Get)

	hHandler := homeHandler.NewHome()
	app.Get("/", hHandler.Get)

	fHandler := filmsHandler.NewFilm(repo, notifier, maxListLength)
	app.Get("/films", fHandler.List)
	app.Post("/films", fHandler.Create)
	app.Delete("/films", fHandler.Delete)
	app.Post("/generators/films/:num<min(5);max(50)>", fHandler.Generate)

	dHandler := filmsHandler.NewDirector(repo, maxListLength)
	app.Get("/directors", dHandler.List)

	// Or extend your config for customization
	app.Use(filesystem.New(filesystem.Config{
		Root: http.FS(assetsFS),
		//PathPrefix:   "/assets",
		Browse:       false,
		NotFoundFile: "404.html",
		MaxAge:       3600,
	}))

	app.Listen(":8000")
}
