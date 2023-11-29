package model

import "github.com/go-playground/validator/v10"

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type Film struct {
	Title    string `validate:"required" fake:"{moviename}"`
	Director string `validate:"required" fake:"{name}"`
}

func (f Film) Clone() Film {
	return Film{Title: f.Title, Director: f.Director}
}

func (f *Film) Validate() error {
	return validate.Struct(f)
}

type Director struct {
	Name   string   `validate:"required"`
	Titles []string `validate:"required"`
}

func (d Director) Clone() Director {
	titlesClone := make([]string, 0, len(d.Titles))
	for _, title := range d.Titles {
		titlesClone = append(titlesClone, title)
	}

	return Director{Name: d.Name, Titles: titlesClone}
}
