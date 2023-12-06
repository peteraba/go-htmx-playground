package view

import (
	"embed"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofiber/template/html/v2"
)

//go:embed templates/*.html
var viewsFS embed.FS

func NewEngine() *html.Engine {
	engine := html.NewFileSystem(http.FS(viewsFS), ".html")
	engine.
		AddFunc("str_slice", func(s []string) template.HTML {
			//nolint: gosec
			return template.HTML(strings.Join(s, ","))
		}).
		AddFunc("url", func(str string) template.URL {
			_, err := url.ParseRequestURI(str)
			if err != nil {
				panic(err)
			}
			//nolint: gosec
			return template.URL(str)
		})

	return engine
}
