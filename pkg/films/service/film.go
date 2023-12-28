package service

import (
	"fmt"
	"log/slog"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/peteraba/go-htmx-playground/lib/pagination"
	"github.com/peteraba/go-htmx-playground/pkg/films/model"
	"github.com/peteraba/go-htmx-playground/pkg/films/repository"
)

type Film struct {
	repo   *repository.FilmRepo
	logger *slog.Logger
}

func NewFilm(repo *repository.FilmRepo, logger *slog.Logger) Film {
	return Film{
		repo:   repo,
		logger: logger,
	}
}

func (f Film) Insert(films ...model.Film) error {
	f.repo.Insert(films...)

	return nil
}

func (f Film) Generate(num int) (int, error) {
	prevCount := f.repo.CountFilms("")

	for i := 0; i < num; i++ {
		//nolint: exhaustruct
		newFilm := model.Film{}

		err := gofakeit.Struct(&newFilm)
		if err != nil {
			return 0, fmt.Errorf("failed to generate new film. num: %d, err: %w", i+1, err)
		}

		f.repo.Insert(newFilm)
	}

	newCount := f.repo.CountFilms("")

	return newCount - prevCount, nil
}

func (f Film) DeleteByTitle(titles ...string) error {
	f.repo.DeleteByTitle(titles...)

	return nil
}

func (f Film) Count(searchTerm string) (int, error) {
	return f.repo.CountFilms(searchTerm), nil
}

func (f Film) Truncate() error {
	_ = f.repo.Truncate()

	return nil
}

func (f Film) List(currentPage, pageSize int, basePath, searchTerm string) ([]model.Film, pagination.Pagination, error) {
	offset := (currentPage - 0) * pageSize

	films, err := f.repo.ListFilms(offset, pageSize, searchTerm)
	if err != nil {
		return nil, pagination.Pagination{}, fmt.Errorf("failed to list films. err: %w", err)
	}

	filmCount := f.repo.CountFilms(searchTerm)

	params := map[string]string{"q": searchTerm}

	filmPagination := pagination.New(currentPage, pageSize, filmCount, basePath, params, "#movie-list")

	return films, filmPagination, nil
}
