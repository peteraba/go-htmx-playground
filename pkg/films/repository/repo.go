package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"sync"

	"github.com/samber/lo"

	"github.com/peteraba/go-htmx-playground/pkg/films/model"
)

type FilmRepo struct {
	lock   sync.RWMutex
	logger *slog.Logger
	// Key is the title of the film
	films      map[string]model.Film
	filmTitles []string
	// Key is the name of the director
	directors     map[string]model.Director
	directorNames []string
	maxLimit      int
}

func (r *FilmRepo) Insert(newFilms ...model.Film) *FilmRepo {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, newFilm := range newFilms {
		// title already exists
		if _, ok := r.films[newFilm.Title]; ok {
			continue
		}

		// set film
		r.films[newFilm.Title] = newFilm.Clone()
		r.filmTitles = append(r.filmTitles, newFilm.Title)

		// set director
		var newDirector model.Director

		if d, ok := r.directors[newFilm.Director]; ok {
			newDirector = d.Clone()
			newDirector.Titles = append(newDirector.Titles, newFilm.Title)
		} else {
			newDirector = model.Director{Name: newFilm.Director, Titles: []string{newFilm.Title}}
			r.directorNames = append(r.directorNames, newFilm.Director)
		}

		sort.Strings(newDirector.Titles)

		r.directors[newFilm.Director] = newDirector
	}

	// reindex
	sort.Strings(r.filmTitles)
	sort.Strings(r.directorNames)

	return r
}

func (r *FilmRepo) deleteOneByTitle(title string) *FilmRepo {
	film, ok := r.films[title]
	if !ok {
		return r
	}

	delete(r.films, title)

	director, ok := r.directors[film.Director]
	if !ok {
		panic(fmt.Sprintf("failed to find director %s", film.Director))
	}

	// Remove title from the director's titles
	r.directors[film.Director] = model.Director{
		Name:   film.Director,
		Titles: lo.Without(director.Titles, title),
	}

	// Remove the director if they do not have any titles
	if len(r.directors[director.Name].Titles) == 0 {
		r.directorNames = lo.Without(r.directorNames, director.Name)
		if r.directorNames == nil {
			r.directorNames = []string{}
		}

		delete(r.directors, director.Name)
	}

	r.filmTitles = lo.Without(r.filmTitles, title)
	if r.filmTitles == nil {
		r.filmTitles = []string{}
	}

	return r
}

func (r *FilmRepo) DeleteByTitle(titles ...string) *FilmRepo {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, title := range titles {
		r.deleteOneByTitle(title)
	}

	return r
}

func (r *FilmRepo) Truncate() *FilmRepo {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.films = make(map[string]model.Film)
	r.filmTitles = []string{}
	r.directors = make(map[string]model.Director)
	r.directorNames = []string{}

	return r
}

var errFilmNotFound = errors.New("film was not found")

// fetchFilmsByTitles is not safe to call concurrently!
func (r *FilmRepo) fetchFilmsByTitles(titles []string) []model.Film {
	films := make([]model.Film, 0, len(titles))

	for _, title := range titles {
		if f, ok := r.films[title]; ok {
			films = append(films, f.Clone())
		} else {
			panic(fmt.Errorf("title: %s, err: %w", title, errFilmNotFound))
		}
	}

	return films
}

var (
	errLimitTooLarge = errors.New("limit too large")
	errLimitTooSmall = errors.New("limit too small")
)

func (r *FilmRepo) ListFilms(offset, limit int) ([]model.Film, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if limit > r.maxLimit {
		return nil, fmt.Errorf("limit > maxLimit. limit: %d, maxLimit: %d, err: %w", limit, r.maxLimit, errLimitTooLarge)
	}

	if limit < 1 {
		return nil, fmt.Errorf("limit < minLimit. limit: %d, minLimit: %d, err: %w", limit, 1, errLimitTooSmall)
	}

	if offset >= len(r.films) {
		return nil, nil
	}

	if offset+limit >= len(r.films) {
		return r.fetchFilmsByTitles(r.filmTitles[offset:]), nil
	}

	return r.fetchFilmsByTitles(r.filmTitles[offset : offset+limit]), nil
}

func (r *FilmRepo) CountFilms() int {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return len(r.films)
}

var errDirectorNotFound = errors.New("director was not found")

// fetchDirectorsByNames is not safe to call concurrently!
func (r *FilmRepo) fetchDirectorsByNames(names []string) []model.Director {
	directors := make([]model.Director, 0, len(names))

	for _, name := range names {
		if d, ok := r.directors[name]; ok {
			directors = append(directors, d.Clone())
		} else {
			r.logger.
				With("director", name).
				With("names", names).
				With("directors", r.directors).
				With("director names", r.directorNames).
				Error("director was not found")

			panic(fmt.Errorf("name: %s, err: %w", name, errDirectorNotFound))
		}
	}

	return directors
}

func (r *FilmRepo) ListDirectors(offset, limit int) ([]model.Director, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if limit > r.maxLimit {
		return nil, fmt.Errorf("limit > maxLimit. limit: %d, maxLimit: %d, err: %w", limit, r.maxLimit, errLimitTooLarge)
	}

	if limit < 1 {
		return nil, fmt.Errorf("limit < minLimit. limit: %d, minLimit: %d, err: %w", limit, 1, errLimitTooSmall)
	}

	if offset >= len(r.films) {
		return nil, nil
	}

	if offset+limit >= len(r.films) {
		return r.fetchDirectorsByNames(r.directorNames[offset:]), nil
	}

	return r.fetchDirectorsByNames(r.directorNames[offset : offset+limit]), nil
}

func (r *FilmRepo) ListAllDirectorNames() []string {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.directorNames
}

func (r *FilmRepo) ListAllTitles() []string {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.filmTitles
}

func (r *FilmRepo) CountDirectors() int {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return len(r.directors)
}

func NewFilmRepo(logger *slog.Logger, maxLimit int, films ...model.Film) *FilmRepo {
	repo := &FilmRepo{
		lock:          sync.RWMutex{},
		logger:        logger,
		filmTitles:    []string{},
		films:         make(map[string]model.Film),
		directorNames: []string{},
		directors:     make(map[string]model.Director),
		maxLimit:      maxLimit,
	}

	if len(films) > 0 {
		repo.Insert(films...)
	}

	return repo
}
