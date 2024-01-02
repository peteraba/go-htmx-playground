package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"sync"

	"github.com/gosimple/slug"
	"github.com/samber/lo"

	"github.com/peteraba/go-htmx-playground/pkg/movies/model"
)

type MovieRepo struct {
	lock   sync.RWMutex
	logger *slog.Logger

	// Key is the slugified keys of the movie
	movies    map[string]model.Movie
	movieKeys []string

	// Key is the slugified name of the director
	directors    map[string]model.Director
	directorKeys []string

	maxLimit int
}

func (r *MovieRepo) Insert(newMovies ...model.Movie) *MovieRepo {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, newMovie := range newMovies {
		movieKey := slug.Make(newMovie.Title)

		// keys already exists
		if _, ok := r.movies[movieKey]; ok {
			continue
		}

		// set movie
		r.movies[movieKey] = newMovie.Clone()
		r.movieKeys = append(r.movieKeys, movieKey)

		// set director
		var newDirector model.Director

		directorKey := slug.Make(newMovie.Director)

		if d, ok := r.directors[directorKey]; ok {
			newDirector = d.Clone()
			newDirector.Titles = append(newDirector.Titles, newMovie.Title)
			sort.Strings(newDirector.Titles)
		} else {
			newDirector = model.Director{ID: directorKey, Name: newMovie.Director, Titles: []string{newMovie.Title}}
			r.directorKeys = append(r.directorKeys, directorKey)
		}

		r.directors[directorKey] = newDirector
	}

	// reindex
	sort.Strings(r.movieKeys)
	sort.Strings(r.directorKeys)

	return r
}

// deleteMovieByKey is not safe to call concurrently!
func (r *MovieRepo) deleteMovieByKey(movieKey string) *MovieRepo {
	movie, ok := r.movies[movieKey]
	if !ok {
		return r
	}

	directorKey := slug.Make(movie.Director)

	director, ok := r.directors[directorKey]
	if !ok {
		panic(fmt.Sprintf("failed to find director %s", movie.Director))
	}

	delete(r.movies, movieKey)

	r.movieKeys = lo.Without(r.movieKeys, movieKey)
	if r.movieKeys == nil {
		r.movieKeys = []string{}
	}

	if len(director.Titles) > 1 {
		r.directors[directorKey] = model.Director{
			ID:     directorKey,
			Name:   movie.Director,
			Titles: lo.Without(director.Titles, movie.Title),
		}
	} else {
		r.directorKeys = lo.Without(r.directorKeys, directorKey)
		if r.directorKeys == nil {
			r.directorKeys = []string{}
		}

		delete(r.directors, directorKey)
	}

	return r
}

func (r *MovieRepo) DeleteMoviesByKey(keys ...string) *MovieRepo {
	r.lock.Lock()
	defer r.lock.Unlock()

	for _, key := range keys {
		r.deleteMovieByKey(key)
	}

	return r
}

func (r *MovieRepo) Truncate() *MovieRepo {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.movies = make(map[string]model.Movie)
	r.movieKeys = []string{}
	r.directors = make(map[string]model.Director)
	r.directorKeys = []string{}

	return r
}

var errMovieNotFound = errors.New("movie was not found")

// fetchMoviesByKeys is not safe to call concurrently!
func (r *MovieRepo) fetchMoviesByKeys(keys []string) []model.Movie {
	movies := make([]model.Movie, 0, len(keys))

	for _, key := range keys {
		if f, ok := r.movies[key]; ok {
			movies = append(movies, f.Clone())
		} else {
			panic(fmt.Errorf("key: %s, err: %w", key, errMovieNotFound))
		}
	}

	return movies
}

var (
	errLimitTooLarge = errors.New("limit too large")
	errLimitTooSmall = errors.New("limit too small")
)

func (r *MovieRepo) ListMovies(offset, limit int, searchTerm string) ([]model.Movie, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if limit > r.maxLimit {
		return nil, fmt.Errorf("limit > maxLimit. limit: %d, maxLimit: %d, err: %w", limit, r.maxLimit, errLimitTooLarge)
	}

	if limit < 1 {
		return nil, fmt.Errorf("limit < minLimit. limit: %d, minLimit: %d, err: %w", limit, 1, errLimitTooSmall)
	}

	if offset >= len(r.movies) {
		return []model.Movie{}, nil
	}

	movies := []model.Movie{}

	for _, key := range r.movieKeys {
		movie, ok := r.movies[key]
		if !ok {
			panic(fmt.Sprintf("failed to find movie for found key: %s", key))
		}

		if searchTerm == "" || strings.Contains(movie.Title, searchTerm) {
			movies = append(movies, movie.Clone())
		}
	}

	if offset > 0 {
		if offset > len(movies) {
			return nil, nil
		}

		movies = movies[offset:]
	}

	if limit < len(movies) {
		movies = movies[:limit]
	}

	return movies, nil
}

func (r *MovieRepo) CountMovies(searchTerm string) int {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if searchTerm == "" {
		return len(r.movies)
	}

	searchTerm = strings.ToLower(searchTerm)

	movies := r.fetchMoviesByKeys(r.movieKeys)
	filteredMovies := lo.Filter(movies, func(f model.Movie, idx int) bool {
		return strings.Contains(strings.ToLower(f.Title), searchTerm)
	})

	return len(filteredMovies)
}

var errDirectorNotFound = errors.New("director was not found")

// fetchDirectorsByKeys is not safe to call concurrently!
func (r *MovieRepo) fetchDirectorsByKeys(keys []string) []model.Director {
	directors := make([]model.Director, 0, len(keys))

	for _, key := range keys {
		if d, ok := r.directors[key]; ok {
			directors = append(directors, d.Clone())
		} else {
			r.logger.
				With("director key searched", key).
				With("keys", keys).
				With("directors", r.directors).
				With("director keys", r.directorKeys).
				Error("director was not found")

			panic(fmt.Errorf("key: %s, err: %w", key, errDirectorNotFound))
		}
	}

	return directors
}

func (r *MovieRepo) ListDirectors(offset, limit int) ([]model.Director, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if limit > r.maxLimit {
		return nil, fmt.Errorf("limit > maxLimit. limit: %d, maxLimit: %d, err: %w", limit, r.maxLimit, errLimitTooLarge)
	}

	if limit < 1 {
		return nil, fmt.Errorf("limit < minLimit. limit: %d, minLimit: %d, err: %w", limit, 1, errLimitTooSmall)
	}

	if offset >= len(r.movies) {
		return nil, nil
	}

	if offset+limit >= len(r.movies) {
		return r.fetchDirectorsByKeys(r.directorKeys[offset:]), nil
	}

	return r.fetchDirectorsByKeys(r.directorKeys[offset : offset+limit]), nil
}

func (r *MovieRepo) ListAllDirectorNames() []string {
	r.lock.RLock()
	defer r.lock.RUnlock()

	keys := lo.MapToSlice(r.directors, func(key string, d model.Director) string { return d.Name })

	sort.Strings(keys)

	return keys
}

func (r *MovieRepo) ListAllTitles() []string {
	r.lock.RLock()
	defer r.lock.RUnlock()

	keys := lo.MapToSlice(r.movies, func(key string, d model.Movie) string { return d.Title })

	sort.Strings(keys)

	return keys
}

func (r *MovieRepo) CountDirectors() int {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return len(r.directors)
}

func NewMovieRepo(logger *slog.Logger, maxLimit int, movies ...model.Movie) *MovieRepo {
	repo := &MovieRepo{
		lock:         sync.RWMutex{},
		logger:       logger,
		movieKeys:    []string{},
		movies:       make(map[string]model.Movie),
		directorKeys: []string{},
		directors:    make(map[string]model.Director),
		maxLimit:     maxLimit,
	}

	if len(movies) > 0 {
		repo.Insert(movies...)
	}

	return repo
}
