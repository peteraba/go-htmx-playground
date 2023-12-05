package view

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/gofiber/template/html/v2"
)

//go:embed templates/*.html
var viewsFS embed.FS

func NewEngine() *html.Engine {
	engine := html.NewFileSystem(http.FS(viewsFS), ".html")
	engine.
		AddFunc("str_slice", func(s []string) template.HTML {
			return template.HTML(strings.Join(s, ","))
		}).
		AddFunc("url", func(s string) template.URL {
			return template.URL(s)
		})

	return engine
}
