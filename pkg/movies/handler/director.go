package handler

import (
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/lib/htmx"
	"github.com/peteraba/go-htmx-playground/lib/jason"
	"github.com/peteraba/go-htmx-playground/lib/log"
	"github.com/peteraba/go-htmx-playground/lib/pagination"
	"github.com/peteraba/go-htmx-playground/pkg/movies/repository"
	"github.com/peteraba/go-htmx-playground/pkg/movies/view"
	notificationsService "github.com/peteraba/go-htmx-playground/pkg/notifications/service"
)

type Director struct {
	repo     *repository.MovieRepo
	pageSize int
	notifier *notificationsService.Notifier
	logger   *slog.Logger
}

func NewDirector(repo *repository.MovieRepo, pageSize int, notifier *notificationsService.Notifier, logger *slog.Logger) Director {
	return Director{
		repo:     repo,
		pageSize: pageSize,
		notifier: notifier,
		logger:   logger,
	}
}

type Response struct {
	Query map[string]interface{} `json:"query,omitempty"`
	Self  string                 `json:"self"`
	First string                 `json:"first,omitempty"`
	Prev  string                 `json:"prev,omitempty"`
	Next  string                 `json:"next,omitempty"`
	Last  string                 `json:"last,omitempty"`
	Items interface{}            `json:"items"`
}

func (d Director) List(c *fiber.Ctx) error {
	currentPage := c.QueryInt("page", 1)
	if currentPage <= 0 {
		currentPage = 1
	}

	offset := (currentPage - 1) * d.pageSize

	directors, err := d.repo.ListDirectors(offset, d.pageSize)
	if err != nil {
		d.logger.Error("Error while listing directors.", log.Err(err))
		d.notifier.Error(err.Error(), c.IP())

		return c.SendStatus(http.StatusInternalServerError)
	}

	p := pagination.New(currentPage, d.pageSize, d.repo.CountDirectors(), c.Path(), nil, "#wrapper")

	if htmx.AcceptHTML(c.GetReqHeaders()) {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

		component := view.DirectorsPage(directors, p.Template())

		return component.Render(c.Context(), c.Response().BodyWriter())
	}

	return jason.SendList(c, directors, p)
}
