package handler

import (
	"net/http"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/lib/htmx"
	"github.com/peteraba/go-htmx-playground/lib/pagination"
	"github.com/peteraba/go-htmx-playground/pkg/films/model"
	"github.com/peteraba/go-htmx-playground/pkg/films/repository"
)

type Film struct {
	repo     *repository.FilmRepo
	pageSize int
}

func NewFilm(repo *repository.FilmRepo, pageSize int) Film {
	return Film{
		repo:     repo,
		pageSize: pageSize,
	}
}

func (f Film) List(c *fiber.Ctx) error {
	return f.list(c)
}

func (f Film) Create(c *fiber.Ctx) error {
	newFilm := model.Film{
		Title:    c.FormValue("title"),
		Director: c.FormValue("director"),
	}
	if newFilm.Validate() != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	f.repo.Insert(newFilm)

	return f.list(c)
}

func (f Film) Generate(c *fiber.Ctx) error {
	n, err := c.ParamsInt("num")
	if err != nil || n < 5 || n >= 50 || !htmx.IsHx(c.GetReqHeaders()) {
		return c.SendStatus(http.StatusBadRequest)
	}

	for i := 0; i < n; i++ {
		newFilm := model.Film{}
		err = gofakeit.Struct(&newFilm)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		f.repo.Insert(newFilm)
	}

	return f.list(c)
}

func (f Film) Delete(c *fiber.Ctx) error {
	f.repo.Truncate()

	return f.list(c)
}

func (f Film) list(c *fiber.Ctx) error {
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

	p := pagination.New(currentPage, f.pageSize, f.repo.CountFilms(), c.Path())

	if htmx.IsHx(c.GetReqHeaders()) {
		return c.Render("views/films", fiber.Map{"Path": c.Path(), "Films": films, "Pagination": p})
	}

	// Render index
	return c.Render("views/films", fiber.Map{"Path": c.Path(), "Films": films, "Pagination": p}, "views/layout")
}
