package films

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	filmsHandler "github.com/peteraba/go-htmx-playground/pkg/films/handler"
	"github.com/peteraba/go-htmx-playground/pkg/films/repository"
	filmsService "github.com/peteraba/go-htmx-playground/pkg/films/service"
	notificationsService "github.com/peteraba/go-htmx-playground/pkg/notifications/service"
)

func Setup(app *fiber.App, logger *slog.Logger, notifier *notificationsService.Notifier, maxListLength int, version string) {
	repo := repository.NewFilmRepo(logger, maxListLength)

	fService := filmsService.NewFilm(repo, logger)
	fHandler := filmsHandler.NewFilm(fService, maxListLength, notifier, logger, version)
	app.Get("/films", fHandler.List)
	app.Post("/films", fHandler.Create)
	app.Post("/films-delete", fHandler.DeleteForm)
	app.Delete("/films", fHandler.Delete)
	app.Post("/generators/films/:num<min(5);max(50)>", fHandler.Generate)

	dHandler := filmsHandler.NewDirector(repo, maxListLength, notifier, logger)
	app.Get("/directors", dHandler.List)
}
