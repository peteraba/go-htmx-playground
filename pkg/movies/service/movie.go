package service

import (
	"fmt"
	"log/slog"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/peteraba/go-htmx-playground/lib/pagination"
	"github.com/peteraba/go-htmx-playground/pkg/movies/model"
	"github.com/peteraba/go-htmx-playground/pkg/movies/repository"
)

type Movie struct {
	repo   *repository.MovieRepo
	logger *slog.Logger
}

func NewMovie(repo *repository.MovieRepo, logger *slog.Logger) *Movie {
	return &Movie{
		repo:   repo,
		logger: logger,
	}
}

func (f Movie) Insert(movies ...model.Movie) error {
	f.repo.Insert(movies...)

	return nil
}

func (f Movie) Generate(num int) (int, error) {
	prevCount := f.repo.CountMovies("")

	for i := 0; i < num; i++ {
		//nolint: exhaustruct
		newMovie := model.Movie{}

		err := gofakeit.Struct(&newMovie)
		if err != nil {
			return 0, fmt.Errorf("failed to generate new movie. num: %d, err: %w", i+1, err)
		}

		f.repo.Insert(newMovie)
	}

	newCount := f.repo.CountMovies("")

	return newCount - prevCount, nil
}

func (f Movie) DeleteByTitle(titles ...string) (int, error) {
	oldCount := f.repo.CountMovies("")

	f.repo.DeleteMoviesByKey(titles...)

	newCount := f.repo.CountMovies("")

	return oldCount - newCount, nil
}

func (f Movie) Count(searchTerm string) (int, error) {
	return f.repo.CountMovies(searchTerm), nil
}

func (f Movie) Truncate() (int, error) {
	oldCount := f.repo.CountMovies("")

	_ = f.repo.Truncate()

	newCount := f.repo.CountMovies("")

	return oldCount - newCount, nil
}

func (f Movie) List(currentPage, pageSize int, basePath, searchTerm string) ([]model.Movie, pagination.Pagination, error) {
	offset := (currentPage - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	movies, err := f.repo.ListMovies(offset, pageSize, searchTerm)
	if err != nil {
		return nil, pagination.Pagination{}, fmt.Errorf("failed to list movies. err: %w", err)
	}

	movieCount := f.repo.CountMovies(searchTerm)

	params := make(map[string]string)
	if searchTerm != "" {
		params["q"] = searchTerm
	}

	p := pagination.New(currentPage, pageSize, movieCount, basePath, params, "#movie-list")

	return movies, p, nil
}
