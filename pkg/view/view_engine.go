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
		AddFunc("str_slice", func(str []string) template.HTML {
			for i, v := range str {
				str[i] = template.HTMLEscapeString(v)
			}

			// nolint: gosec
			return template.HTML(strings.Join(str, ", "))
		}).
		AddFunc("url", func(str string) template.URL {
			_, err := url.ParseRequestURI(str)
			if err != nil {
				return ""
			}

			// nolint: gosec
			return template.URL(str)
		})

	return engine
}
