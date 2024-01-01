package movies

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/pkg/movies/handler"
	"github.com/peteraba/go-htmx-playground/pkg/movies/repository"
	"github.com/peteraba/go-htmx-playground/pkg/movies/service"
	notificationsService "github.com/peteraba/go-htmx-playground/pkg/notifications/service"
)

func Setup(app *fiber.App, logger *slog.Logger, notifier *notificationsService.Notifier, maxListLength int, version string) {
	repo := repository.NewMovieRepo(logger, maxListLength)

	fService := service.NewMovie(repo, logger)
	fHandler := handler.NewMovie(fService, maxListLength, notifier, logger, version)
	app.Get("/movies", fHandler.List)
	app.Post("/movies", fHandler.Create)
	app.Post("/movies-delete", fHandler.DeleteForm)
	app.Delete("/movies", fHandler.Delete)
	app.Post("/generators/movies/:num<min(5);max(50)>", fHandler.Generate)

	dHandler := handler.NewDirector(repo, maxListLength, notifier, logger)
	app.Get("/directors", dHandler.List)
}
