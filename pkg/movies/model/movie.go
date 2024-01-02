package model

import (
	"github.com/go-playground/validator/v10"
)

// nolint: gochecknoglobals
var validate *validator.Validate

// nolint: gochecknoinits
func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type Movie struct {
	ID       string `fake:"{uuid}"      json:"id"`
	Title    string `fake:"{moviename}" json:"title"    validate:"required"`
	Director string `fake:"{name}"      json:"director" validate:"required"`
}

func (f Movie) Clone() Movie {
	return Movie{ID: f.ID, Title: f.Title, Director: f.Director}
}

func (f Movie) Validate() error {
	//nolint: wrapcheck
	return validate.Struct(&f)
}

type Director struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"   validate:"required"`
	Titles []string `json:"titles" validate:"required"`
}

func (d Director) Clone() Director {
	titlesClone := make([]string, 0, len(d.Titles))
	titlesClone = append(titlesClone, d.Titles...)

	return Director{ID: d.ID, Name: d.Name, Titles: titlesClone}
}
