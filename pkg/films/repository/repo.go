package repository

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/peteraba/go-htmx-playground/pkg/films/model"
)

type FilmRepo struct {
	lock          sync.RWMutex
	films         map[string]model.Film
	titles        []string
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
		r.titles = append(r.titles, newFilm.Title)

		// set director
		newDirector := model.Director{Name: newFilm.Director, Titles: []string{newFilm.Title}}
		if d, ok := r.directors[newFilm.Director]; ok {
			newDirector.Titles = append(d.Titles, newFilm.Title)
		} else {
			r.directorNames = append(r.directorNames, newFilm.Director)
		}
		r.directors[newFilm.Director] = newDirector
	}

	// reindex
	sort.Strings(r.titles)
	sort.Strings(r.directorNames)

	return r
}

func (r *FilmRepo) Truncate() *FilmRepo {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.films = make(map[string]model.Film)
	r.titles = []string{}
	r.directors = make(map[string]model.Director)
	r.directorNames = []string{}

	return r
}

// fetchFilmsByTitles is not safe to call concurrently!
func (r *FilmRepo) fetchFilmsByTitles(titles []string) []model.Film {
	films := make([]model.Film, 0, len(titles))

	for _, title := range titles {
		if f, ok := r.films[title]; ok {
			films = append(films, f.Clone())
		} else {
			fmt.Printf("%v", r.films)
			panic(fmt.Errorf("film was not found: %s", title))
		}
	}

	return films
}

func (r *FilmRepo) ListFilms(offset, limit int) ([]model.Film, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if limit > r.maxLimit {
		return nil, fmt.Errorf("limit too large: %d (max: %d)", limit, r.maxLimit)
	}
	if limit < 1 {
		return nil, fmt.Errorf("invalid too small: %d (min: 1)", limit)
	}

	if offset >= len(r.films) {
		return nil, nil
	}

	if offset+limit >= len(r.films) {
		return r.fetchFilmsByTitles(r.titles[offset:]), nil
	}

	return r.fetchFilmsByTitles(r.titles[offset : offset+limit]), nil
}

func (r *FilmRepo) CountFilms() int {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return len(r.films)
}

// fetchDirectorsByNames is not safe to call concurrently!
func (r *FilmRepo) fetchDirectorsByNames(names []string) []model.Director {
	directors := make([]model.Director, 0, len(names))

	for _, name := range names {
		if d, ok := r.directors[name]; ok {
			directors = append(directors, d.Clone())
		} else {
			log.Printf("names: %v", names)
			log.Printf("directors: %v", r.directors)
			log.Printf("director names: %v", r.directorNames)
			panic(fmt.Errorf("director was not found: %s", name))
		}
	}

	return directors
}

func (r *FilmRepo) ListDirectors(offset, limit int) ([]model.Director, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	if limit > r.maxLimit {
		return nil, fmt.Errorf("limit too large: %d (max: %d)", limit, r.maxLimit)
	}
	if limit < 1 {
		return nil, fmt.Errorf("invalid too small: %d (min: 1)", limit)
	}

	if offset >= len(r.films) {
		return nil, nil
	}

	if offset+limit >= len(r.films) {
		return r.fetchDirectorsByNames(r.directorNames[offset:]), nil
	}

	return r.fetchDirectorsByNames(r.directorNames[offset : offset+limit]), nil
}

func (r *FilmRepo) CountDirectors() int {
	r.lock.RLock()
	defer r.lock.RUnlock()

	return len(r.directors)
}

func NewFilmRepo(maxLimit int, films ...model.Film) *FilmRepo {
	repo := &FilmRepo{
		films:     make(map[string]model.Film),
		directors: make(map[string]model.Director),
		maxLimit:  maxLimit,
	}

	if len(films) > 0 {
		repo.Insert(films...)
	}

	return repo
}
