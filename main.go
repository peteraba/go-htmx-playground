package main

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"

	"github.com/peteraba/go-htmx-playground/lib/htmx"
	"github.com/peteraba/go-htmx-playground/pkg/films/handler"
	"github.com/peteraba/go-htmx-playground/pkg/films/repository"
)

//go:embed views/*.html
var viewsFS embed.FS

//go:embed assets/*
var assetsFS embed.FS

var themes = []string{
	"light",
	"dark",
	"cupcake",
	"bumblebee",
	"emerald",
	"corporate",
	"synthwave",
	"retro",
	"cyberpunk",
	"valentine",
	"halloween",
	"garden",
	"forest",
	"aqua",
	"lofi",
	"pastel",
	"fantasy",
	"wireframe",
	"black",
	"luxury",
	"dracula",
	"cmyk",
	"autumn",
	"business",
	"acid",
	"lemonade",
	"night",
	"coffee",
	"winter",
	"dim",
	"nord",
	"sunset",
}

func main() {
	const maxListLength = 10

	engine := html.NewFileSystem(http.FS(viewsFS), ".html")
	engine.AddFunc(
		"str_slice", func(s []string) template.HTML {
			return template.HTML(strings.Join(s, ","))
		},
	)

	app := fiber.New(fiber.Config{
		Views:     engine,
		Immutable: true,
	})

	repo := repository.NewFilmRepo(maxListLength)

	app.Get("/", func(c *fiber.Ctx) error {
		if htmx.IsHx(c.GetReqHeaders()) {
			return c.Render("views/home", fiber.Map{"Path": c.Path()})
		}
		// Render index
		return c.Render("views/home", fiber.Map{"Path": c.Path()}, "views/layout")
	})

	app.Get("/colors", func(c *fiber.Ctx) error {
		if htmx.IsHx(c.GetReqHeaders()) {
			return c.Render("views/colors", fiber.Map{"Path": c.Path(), "Themes": themes})
		}
		// Render index
		return c.Render("views/colors", fiber.Map{"Path": c.Path(), "Themes": themes}, "views/layout")
	})

	filmHandler := handler.NewFilm(repo, maxListLength)
	app.Get("/films", filmHandler.List)
	app.Post("/films", filmHandler.Create)
	app.Delete("/films", filmHandler.Delete)
	app.Post("/generators/films/:num<min(5);max(50)>", filmHandler.Generate)

	directorHandler := handler.NewDirector(repo, maxListLength)
	app.Get("/directors", directorHandler.List)

	app.Get("/assets/:file", func(c *fiber.Ctx) error {
		fileName := "assets/" + c.Params("file")

		content, err := assetsFS.ReadFile(fileName)
		if err != nil {
			return err
		}

		return c.Send(content)
	})

	app.Listen(":8000")
}
