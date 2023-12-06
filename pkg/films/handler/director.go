package handler

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/lib/htmx"
	"github.com/peteraba/go-htmx-playground/lib/pagination"
	"github.com/peteraba/go-htmx-playground/pkg/films/repository"
)

type Director struct {
	repo     *repository.FilmRepo
	pageSize int
}

func NewDirector(repo *repository.FilmRepo, pageSize int) Director {
	return Director{
		repo:     repo,
		pageSize: pageSize,
	}
}

func (d Director) List(c *fiber.Ctx) error {
	bind := fiber.Map{"Path": c.Path(), "Url": c.BaseURL()}

	currentPage := c.QueryInt("page", 1)
	if currentPage <= 0 {
		currentPage = 1
	}

	offset := (currentPage - 1) * d.pageSize

	directors, err := d.repo.ListDirectors(offset, d.pageSize)
	if err != nil {
		c.Status(http.StatusInternalServerError)

		return fmt.Errorf("failed to list directors, err: %w", err)
	}

	bind["Directors"] = directors
	bind["Pagination"] = pagination.New(currentPage, d.pageSize, d.repo.CountDirectors(), c.Path())

	if htmx.IsHx(c.GetReqHeaders()) {
		return c.Render("templates/directors", bind)
	}

	return c.Render("templates/directors", bind, "templates/layout")
}
