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
	Title    string `fake:"{moviename}" validate:"required"`
	Director string `fake:"{name}"      validate:"required"`
}

func (f Movie) Clone() Movie {
	return Movie{Title: f.Title, Director: f.Director}
}

func (f Movie) Validate() error {
	//nolint: wrapcheck
	return validate.Struct(&f)
}

type Director struct {
	Name   string   `validate:"required"`
	Titles []string `validate:"required"`
}

func (d Director) Clone() Director {
	titlesClone := make([]string, 0, len(d.Titles))
	titlesClone = append(titlesClone, d.Titles...)

	return Director{Name: d.Name, Titles: titlesClone}
}
