package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/lib/htmx"
	"github.com/peteraba/go-htmx-playground/lib/jason"
	"github.com/peteraba/go-htmx-playground/lib/log"
	"github.com/peteraba/go-htmx-playground/lib/pagination"
	"github.com/peteraba/go-htmx-playground/pkg/movies/model"
	"github.com/peteraba/go-htmx-playground/pkg/movies/service"
	"github.com/peteraba/go-htmx-playground/pkg/movies/view"
	notificationsService "github.com/peteraba/go-htmx-playground/pkg/notifications/service"
)

type Movie struct {
	service      *service.Movie
	pageSize     int
	notifier     *notificationsService.Notifier
	logger       *slog.Logger
	buildVersion string
}

func NewMovie(movieService *service.Movie, pageSize int, notifier *notificationsService.Notifier, logger *slog.Logger, version string) Movie {
	return Movie{
		service:      movieService,
		pageSize:     pageSize,
		notifier:     notifier,
		logger:       logger,
		buildVersion: version,
	}
}

func (f Movie) getVersion() string {
	if f.buildVersion == "development" {
		return time.Now().Format(time.RFC3339Nano)
	}

	return f.buildVersion
}

func (f Movie) List(c *fiber.Ctx) error {
	return f.list(c, c.Path())
}

func (f Movie) Create(c *fiber.Ctx) error {
	var newMovie model.Movie

	if htmx.AcceptHTML(c.GetReqHeaders()) {
		newMovie = model.Movie{
			Title:    c.FormValue("title"),
			Director: c.FormValue("director"),
		}
	} else if err := c.BodyParser(&newMovie); err != nil {
		f.logger.Error("Error while parsing request body.", log.Err(err))

		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	err := newMovie.Validate()
	if err != nil {
		f.logger.Error("Error while validating new movie.", log.Err(err))
		f.notifier.Error(err.Error(), c.IP())

		return c.SendStatus(http.StatusBadRequest)
	}

	err = f.service.Insert(newMovie)
	if err != nil {
		f.logger.Error("Error inserting the new movie.", log.Err(err))
		f.notifier.Error(err.Error(), c.IP())

		return c.SendStatus(http.StatusInternalServerError)
	}

	count, err := f.service.Count("")
	if err != nil {
		f.logger.Error("Error retrieving the new count.", log.Err(err))
		f.notifier.Error(err.Error(), c.IP())

		return c.SendStatus(http.StatusInternalServerError)
	}

	f.logger.With("title", newMovie.Title, "director", newMovie.Director, "count", count).Info("Added new movie.")
	f.notifier.Success(fmt.Sprintf("`%s` added.", newMovie.Title), c.IP())

	return f.list(c, "/movies")
}

func (f Movie) Generate(c *fiber.Ctx) error {
	randomNumber, err := c.ParamsInt("num")
	if err != nil || randomNumber < 5 || randomNumber >= 50 || !htmx.IsHx(c.GetReqHeaders()) {
		return c.SendStatus(http.StatusBadRequest)
	}

	generated, err := f.service.Generate(randomNumber)
	if err != nil {
		f.logger.Error("Error while generating movies.", log.Err(err))
		f.notifier.Error(err.Error(), c.IP())

		return c.SendStatus(http.StatusBadRequest)
	}

	f.logger.Info(fmt.Sprintf("%d unique movies generated.", generated))
	f.notifier.Success(fmt.Sprintf("%d unique movies generated.", generated), c.IP())

	return f.list(c, "/movies")
}

// DeleteForm is a handler for deleting movies for browsers without JS support enabled.
func (f Movie) DeleteForm(c *fiber.Ctx) error {
	titles, err := f.getMoviesToDelete(c)
	if err != nil {
		return err
	}

	count, err := f.service.DeleteByTitle(titles...)
	if err != nil {
		f.logger.Error("Error while deleting movies.", log.Err(err))
		f.notifier.Error(err.Error(), c.IP())

		return c.SendStatus(http.StatusInternalServerError)
	}

	f.logger.Info(fmt.Sprintf("%d unique movies deleted.", count))
	f.notifier.Info(fmt.Sprintf("%d unique movies deleted.", count), c.IP())

	return c.Redirect("/movies", http.StatusMovedPermanently)
}

// Delete is a handler which handles truncating movies and individual deletes for browsers with JS support enabled.
func (f Movie) Delete(c *fiber.Ctx) error {
	if c.Query("truncate") == "true" || c.Query("truncate") == "1" {
		return f.truncate(c)
	}

	titles, err := f.getMoviesToDelete(c)
	if err != nil {
		f.logger.Error("Error while getting movies to delete.", log.Err(err))
		f.notifier.Error(err.Error(), c.IP())

		return c.SendStatus(http.StatusInternalServerError)
	}

	return f.deleteTitles(c, titles)
}

func (f Movie) deleteTitles(c *fiber.Ctx, titles []string) error {
	count, err := f.service.DeleteByTitle(titles...)
	if err != nil {
		f.logger.Error("Error while deleting movies.", log.Err(err))
		f.notifier.Error(err.Error(), c.IP())

		return c.SendStatus(http.StatusInternalServerError)
	}

	f.logger.Info(fmt.Sprintf("%d unique movies deleted.", count))
	f.notifier.Success(fmt.Sprintf("Movies deleted: %d.", count), c.IP())

	return f.list(c, "/movies")
}

// Delete is a handler which handles truncating movies and individual deletes for browsers with JS support enabled.
func (f Movie) truncate(c *fiber.Ctx) error {
	count, err := f.service.Count("")
	if err != nil {
		f.logger.Error("error counting movies", log.Err(err))
		f.notifier.Info(fmt.Sprintf("Error counting movies: %s", err), c.IP())

		return f.list(c, "/movies")
	}

	_, err = f.service.Truncate()
	if err != nil {
		f.logger.Error("error truncating movies", log.Err(err))
		f.notifier.Error(fmt.Sprintf("Error truncating movies: %s", err), c.IP())

		return c.SendStatus(http.StatusInternalServerError)
	}

	f.logger.Info(fmt.Sprintf("%d unique movies truncated.", count))
	f.notifier.Success(fmt.Sprintf("Movies truncated: %d.", count), c.IP())

	return f.list(c, "/movies")
}

type ExpectedPayload struct {
	Movies []string
}

func (f Movie) getMoviesToDelete(c *fiber.Ctx) ([]string, error) {
	body := new(ExpectedPayload)
	if err := c.BodyParser(body); err != nil {
		return nil, fmt.Errorf("failed to parse request body, err: %w", err)
	}

	return body.Movies, nil
}

func (f Movie) list(c *fiber.Ctx, basePath string) error {
	searchTerm := c.Query("q")

	currentPage := c.QueryInt("page", 1)
	if currentPage <= 0 {
		currentPage = 1
	}

	movies, p, err := f.service.List(currentPage, f.pageSize, basePath, searchTerm)
	if err != nil {
		f.logger.Error("Error while listing movies.", log.Err(err))

		return c.SendStatus(http.StatusInternalServerError)
	}

	if htmx.AcceptHTML(c.GetReqHeaders()) {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

		return f.render(c, movies, p, searchTerm)
	}

	return jason.SendList(c, movies, p)
}

func (f Movie) render(c *fiber.Ctx, movies []model.Movie, moviePagination pagination.Pagination, searchTerm string) error {
	var component templ.Component

	switch htmx.GetTarget(c.GetReqHeaders()) {
	case "movie-list", "#movie-list":
		component = view.MovieList(movies, moviePagination.Template())

	default:
		component = view.MoviesPage(movies, moviePagination.Template(), searchTerm, f.getVersion())
	}

	return component.Render(c.Context(), c.Response().BodyWriter())
}
