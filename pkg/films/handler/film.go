package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/lib/htmx"
	"github.com/peteraba/go-htmx-playground/lib/log"
	"github.com/peteraba/go-htmx-playground/lib/pagination"
	"github.com/peteraba/go-htmx-playground/pkg/films/model"
	"github.com/peteraba/go-htmx-playground/pkg/films/service"
	"github.com/peteraba/go-htmx-playground/pkg/films/view"
	notificationsService "github.com/peteraba/go-htmx-playground/pkg/notifications/service"
)

type Film struct {
	service      *service.Film
	pageSize     int
	notifier     *notificationsService.Notifier
	logger       *slog.Logger
	buildVersion string
}

func NewFilm(filmService *service.Film, pageSize int, notifier *notificationsService.Notifier, logger *slog.Logger, version string) Film {
	return Film{
		service:      filmService,
		pageSize:     pageSize,
		notifier:     notifier,
		logger:       logger,
		buildVersion: version,
	}
}

func (f Film) getVersion() string {
	if f.buildVersion == "development" {
		return time.Now().Format(time.RFC3339Nano)
	}

	return f.buildVersion
}

func (f Film) List(c *fiber.Ctx) error {
	return f.list(c, c.Path())
}

func (f Film) Create(c *fiber.Ctx) error {
	newFilm := model.Film{
		Title:    c.FormValue("title"),
		Director: c.FormValue("director"),
	}

	err := newFilm.Validate()
	if err != nil {
		f.logger.Error("Error while validating new film.", log.Err(err))
		f.notifier.Error(err.Error(), c.IP())

		return c.SendStatus(http.StatusBadRequest)
	}

	f.logger.With("title", newFilm.Title, "director", newFilm.Director).Info("Added new film.")
	f.notifier.Success(fmt.Sprintf("`%s` added.", newFilm.Title), c.IP())

	return f.list(c, "/films")
}

func (f Film) Generate(c *fiber.Ctx) error {
	randomNumber, err := c.ParamsInt("num")
	if err != nil || randomNumber < 5 || randomNumber >= 50 || !htmx.IsHx(c.GetReqHeaders()) {
		return c.SendStatus(http.StatusBadRequest)
	}

	generated, err := f.service.Generate(randomNumber)
	if err != nil {
		f.logger.Error("Error while generating films.", log.Err(err))
		f.notifier.Error(err.Error(), c.IP())

		return c.SendStatus(http.StatusBadRequest)
	}

	f.logger.Info(fmt.Sprintf("%d unique films generated.", generated))
	f.notifier.Success(fmt.Sprintf("%d unique films generated.", generated), c.IP())

	return f.list(c, "/films")
}

// DeleteForm is a handler for deleting films for browsers without JS support enabled.
func (f Film) DeleteForm(c *fiber.Ctx) error {
	titles, err := f.getFilmsToDelete(c)
	if err != nil {
		return err
	}

	count, err := f.service.DeleteByTitle(titles...)
	if err != nil {
		f.logger.Error("Error while deleting films.", log.Err(err))
		f.notifier.Error(err.Error(), c.IP())

		return c.SendStatus(http.StatusInternalServerError)
	}

	f.logger.Info(fmt.Sprintf("%d unique films deleted.", count))
	f.notifier.Info(fmt.Sprintf("%d unique films deleted.", count), c.IP())

	return c.Redirect("/films", http.StatusMovedPermanently)
}

// Delete is a handler which handles truncating films and individual deletes for browsers with JS support enabled.
func (f Film) Delete(c *fiber.Ctx) error {
	if c.Query("truncate") == "true" || c.Query("truncate") == "1" {
		return f.truncate(c)
	}

	titles, err := f.getFilmsToDelete(c)
	if err != nil {
		f.logger.Error("Error while getting films to delete.", log.Err(err))
		f.notifier.Error(err.Error(), c.IP())

		return c.SendStatus(http.StatusInternalServerError)
	}

	return f.deleteTitles(c, titles)
}

func (f Film) deleteTitles(c *fiber.Ctx, titles []string) error {
	count, err := f.service.DeleteByTitle(titles...)
	if err != nil {
		f.logger.Error("Error while deleting films.", log.Err(err))
		f.notifier.Error(err.Error(), c.IP())

		return c.SendStatus(http.StatusInternalServerError)
	}

	f.logger.Info(fmt.Sprintf("%d unique films deleted.", count))
	f.notifier.Success(fmt.Sprintf("Films deleted: %d.", count), c.IP())

	return f.list(c, "/films")
}

// Delete is a handler which handles truncating films and individual deletes for browsers with JS support enabled.
func (f Film) truncate(c *fiber.Ctx) error {
	count, err := f.service.Count("")
	if err != nil {
		f.logger.Error("error counting films", log.Err(err))
		f.notifier.Info(fmt.Sprintf("Error counting films: %s", err), c.IP())

		return f.list(c, "/films")
	}

	_, err = f.service.Truncate()
	if err != nil {
		f.logger.Error("error truncating films", log.Err(err))
		f.notifier.Error(fmt.Sprintf("Error truncating films: %s", err), c.IP())

		return c.SendStatus(http.StatusInternalServerError)
	}

	f.logger.Info(fmt.Sprintf("%d unique films truncated.", count))
	f.notifier.Success(fmt.Sprintf("Films truncated: %d.", count), c.IP())

	return f.list(c, "/films")
}

type ExpectedPayload struct {
	Films []string
}

func (f Film) getFilmsToDelete(c *fiber.Ctx) ([]string, error) {
	body := new(ExpectedPayload)
	if err := c.BodyParser(body); err != nil {
		return nil, fmt.Errorf("failed to parse request body, err: %w", err)
	}

	return body.Films, nil
}

func (f Film) list(c *fiber.Ctx, basePath string) error {
	searchTerm := c.Query("q")

	currentPage := c.QueryInt("page", 1)
	if currentPage <= 0 {
		currentPage = 1
	}

	films, filmPagination, err := f.service.List(currentPage, f.pageSize, basePath, searchTerm)
	if err != nil {
		f.logger.Error("Error while listing films.", log.Err(err))

		return c.SendStatus(http.StatusInternalServerError)
	}

	return f.render(c, films, filmPagination, searchTerm)
}

func (f Film) render(c *fiber.Ctx, films []model.Film, filmPagination pagination.Pagination, searchTerm string) error {
	var component templ.Component

	switch htmx.GetTarget(c.GetReqHeaders()) {
	case "movie-list", "#movie-list":
		component = view.FilmList(films, filmPagination.Template())

	default:
		component = view.FilmsPage(films, filmPagination.Template(), searchTerm, f.getVersion())
	}

	return component.Render(c.Context(), c.Response().BodyWriter())
}
