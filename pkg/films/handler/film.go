package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/lib/htmx"
	"github.com/peteraba/go-htmx-playground/lib/pagination"
	"github.com/peteraba/go-htmx-playground/pkg/films/model"
	"github.com/peteraba/go-htmx-playground/pkg/films/repository"
	notificationsService "github.com/peteraba/go-htmx-playground/pkg/notifications/service"
)

type Film struct {
	logger   *slog.Logger
	repo     *repository.FilmRepo
	notifier *notificationsService.Notifier
	pageSize int
}

func NewFilm(logger *slog.Logger, repo *repository.FilmRepo, notifier *notificationsService.Notifier, pageSize int) Film {
	return Film{
		logger:   logger,
		repo:     repo,
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

	f.repo.Insert(newFilm)
	f.notifier.Success(fmt.Sprintf("`%s` added.", newFilm.Title), c.IP())

	return f.list(c, "/films")
}

func (f Film) Generate(c *fiber.Ctx) error {
	n, err := c.ParamsInt("num")
	if err != nil || n < 5 || n >= 50 || !htmx.IsHx(c.GetReqHeaders()) {
		return c.SendStatus(http.StatusBadRequest)
	}

	prevCount := f.repo.CountFilms()
	for i := 0; i < n; i++ {
		newFilm := model.Film{}
		err = gofakeit.Struct(&newFilm)
		if err != nil {
			f.notifier.Error(err.Error(), c.IP())
			return c.SendStatus(http.StatusInternalServerError)
		}

		f.repo.Insert(newFilm)
	}
	newCount := f.repo.CountFilms()

	f.notifier.Info(fmt.Sprintf("%d unique films generated.", newCount-prevCount), c.IP())

	return f.list(c, "/films")
}

// DeleteForm is a handler for deleting films for browsers without JS support enabled.
func (f Film) DeleteForm(c *fiber.Ctx) error {
	titles, err := f.getFilmsToDelete(c)
	if err != nil {
		return err
	}

	f.repo.DeleteByTitle(titles...)

	return c.Redirect("/titles", http.StatusMovedPermanently)
}

// Delete is a handler which handles truncating films and individual deletes for browsers with JS support enabled.
func (f Film) Delete(c *fiber.Ctx) error {
	titles, _ := f.getFilmsToDelete(c)
	if len(titles) > 0 {
		f.logger.Debug("JS support enabled. Deleting films...")
		f.repo.DeleteByTitle(titles...)
		f.notifier.Success(fmt.Sprintf("Films deleted: %d", len(titles)), c.IP())
	} else if f.repo.CountFilms() == 0 {
		f.logger.Debug("No films to delete.")
		f.notifier.Info("No titles to delete.", c.IP())
	} else {
		f.logger.Debug("JS support not enabled. Truncating films...")
		f.repo.Truncate()
		f.notifier.Success(fmt.Sprintf("All titles deleted."), c.IP())
	}

	return f.list(c, "/titles")
}

type ExpectedPayload struct {
	Films []string
}

func (f Film) getFilmsToDelete(c *fiber.Ctx) ([]string, error) {
	body := new(ExpectedPayload)
	if err := c.BodyParser(body); err != nil {
		return nil, err
	}

	return body.Films, nil
}

func (f Film) list(c *fiber.Ctx, basePath string) error {
	bind := fiber.Map{"Path": c.Path(), "Url": c.BaseURL()}

	currentPage := c.QueryInt("page", 1)
	if currentPage <= 0 {
		currentPage = 1
	}
	offset := (currentPage - 1) * f.pageSize

	films, err := f.repo.ListFilms(offset, f.pageSize)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return err
	}

	bind["Films"] = films
	bind["Pagination"] = pagination.New(currentPage, f.pageSize, f.repo.CountFilms(), basePath)

	if htmx.IsHx(c.GetReqHeaders()) {
		return c.Render("templates/films", bind)
	}

	// Render index
	return c.Render("templates/films", bind, "templates/layout")
}
