package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/lib/htmx"
	"github.com/peteraba/go-htmx-playground/lib/pagination"
	"github.com/peteraba/go-htmx-playground/pkg/films/model"
	"github.com/peteraba/go-htmx-playground/pkg/films/service"
	"github.com/peteraba/go-htmx-playground/pkg/films/view"
	notificationsService "github.com/peteraba/go-htmx-playground/pkg/notifications/service"
)

type Film struct {
	logger   *slog.Logger
	service  service.Film
	pageSize int
	notifier *notificationsService.Notifier
}

func NewFilm(filmService service.Film, pageSize int, notifier *notificationsService.Notifier, logger *slog.Logger) Film {
	return Film{
		logger:   logger,
		service:  filmService,
		pageSize: pageSize,
		notifier: notifier,
	}
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
		f.notifier.Error(err.Error(), c.IP())

		return c.SendStatus(http.StatusBadRequest)
	}

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
		f.notifier.Error(err.Error(), c.IP())

		return c.SendStatus(http.StatusBadRequest)
	}

	f.notifier.Info(fmt.Sprintf("%d unique films generated.", generated), c.IP())

	return f.list(c, "/films")
}

// DeleteForm is a handler for deleting films for browsers without JS support enabled.
func (f Film) DeleteForm(c *fiber.Ctx) error {
	titles, err := f.getFilmsToDelete(c)
	if err != nil {
		return err
	}

	err = f.service.DeleteByTitle(titles...)
	if err != nil {
		f.notifier.Error(err.Error(), c.IP())

		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Redirect("/titles", http.StatusMovedPermanently)
}

// Delete is a handler which handles truncating films and individual deletes for browsers with JS support enabled.
func (f Film) Delete(c *fiber.Ctx) error {
	titles, err := f.getFilmsToDelete(c)
	if err != nil {
		return f.truncate(c)
	}

	return f.deleteTitles(c, titles)
}

func (f Film) deleteTitles(c *fiber.Ctx, titles []string) error {
	f.logger.Debug("JS support enabled. Deleting films...")

	err := f.service.DeleteByTitle(titles...)
	if err != nil {
		f.notifier.Error(err.Error(), c.IP())

		return c.SendStatus(http.StatusInternalServerError)
	}

	f.notifier.Success(fmt.Sprintf("Films deleted: %d", len(titles)), c.IP())

	return f.list(c, "/titles")
}

// Delete is a handler which handles truncating films and individual deletes for browsers with JS support enabled.
func (f Film) truncate(c *fiber.Ctx) error {
	count, _ := f.service.Count("")
	if count == 0 {
		f.logger.Debug("No films to delete.")
		f.notifier.Info("No titles to delete.", c.IP())
	} else {
		f.logger.Debug("JS support not enabled. Truncating films...")
		_ = f.service.Truncate()
	}

	return f.list(c, "/titles")
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
		f.logger.With("err", err).Error("Error while listing films.")

		return c.SendStatus(http.StatusInternalServerError)
	}

	return render(c, films, filmPagination, searchTerm)
}

func render(c *fiber.Ctx, films []model.Film, filmPagination pagination.Pagination, searchTerm string) error {
	var component templ.Component

	switch htmx.GetTarget(c.GetReqHeaders()) {
	case "movie-list", "#movie-list":
		component = view.FilmList(films, filmPagination.Template(), searchTerm)

	default:
		component = view.FilmsPage(films, filmPagination.Template(), searchTerm)
	}

	return component.Render(c.Context(), c.Response().BodyWriter())
}
