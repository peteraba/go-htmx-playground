package main

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"net/http"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

type Film struct {
	Title    string `validate:"required" fake:"{moviename}"`
	Director string `validate:"required" fake:"{name}"`
}

func (f *Film) Validate() error {
	return validate.Struct(f)
}

//go:embed views/*.html
var viewsFS embed.FS

func isHx(headers map[string][]string) bool {
	if headers == nil || len(headers) == 0 {
		return false
	}

	if v, ok := headers["Hx-Request"]; ok && len(v) > 0 {
		return v[0] == "true"
	}

	return false
}

func main() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	engine := html.NewFileSystem(http.FS(viewsFS), ".html")
	engine.AddFunc(
		// add unescape function
		"unescape", func(s string) template.HTML {
			return template.HTML(s)
		},
	)

	app := fiber.New(fiber.Config{Views: engine})

	app.Get("/", func(c *fiber.Ctx) error {
		if isHx(c.GetReqHeaders()) {
			return c.Render("views/home", fiber.Map{"Path": c.Path()})
		}
		// Render index
		return c.Render("views/home", fiber.Map{"Path": c.Path()}, "views/layout")
	})

	app.Get("/films", func(c *fiber.Ctx) error {
		films := []Film{
			{Title: "The Godfather", Director: "Francis Ford Coppola"},
			{Title: "Blade Runner", Director: "Ridley Scott"},
			{Title: "The Thing", Director: "John Carpenter"},
		}

		if isHx(c.GetReqHeaders()) {
			return c.Render("views/films", fiber.Map{"Path": c.Path(), "Films": films})
		}
		// Render index
		return c.Render("views/films", fiber.Map{"Path": c.Path()}, "views/layout")
	})

	app.Get("/directors", func(c *fiber.Ctx) error {
		if isHx(c.GetReqHeaders()) {
			return c.Render("views/directors", fiber.Map{"Path": c.Path()})
		}

		return c.Render("views/directors", fiber.Map{"Path": c.Path()}, "views/layout")
	})

	app.Post("/films", func(c *fiber.Ctx) error {
		f := Film{
			Title:    c.FormValue("title"),
			Director: c.FormValue("director"),
		}
		if f.Validate() != nil {
			return c.SendStatus(http.StatusBadRequest)
		}

		return c.Render("film-list-element", f)
	})

	app.Post("/generators/films/:num<min(5);max(50)>", func(c *fiber.Ctx) error {
		n, err := c.ParamsInt("num")
		if err != nil || n < 5 || n >= 50 {
			return c.SendStatus(http.StatusBadRequest)
		}

		f := Film{}
		b := bytes.Buffer{}
		w := io.Writer(&b)
		tmpl := template.Must(template.ParseFS(viewsFS, "views/films.html"))
		for i := 0; i < n; i++ {
			err = gofakeit.Struct(&f)
			if err != nil {
				return c.SendStatus(http.StatusInternalServerError)
			}

			err = tmpl.ExecuteTemplate(w, "film-list-element", f)
			if err != nil {
				return c.SendStatus(http.StatusInternalServerError)
			}
		}

		return c.Send(b.Bytes())
	})

	app.Listen(":8000")

	// h2 := func(w http.ResponseWriter, r *http.Request) {
	// 	title := r.PostFormValue("title")
	// 	director := r.PostFormValue("director")
	// 	if title == "" || director == "" {
	// 		log.Printf("empty form field. title: %s, director: %s\n", title, director)
	// 		return
	// 	}
	//
	// 	log.Printf("title: %s, director: %s\n", title, director)
	// 	tmpl := template.Must(template.ParseFS(index, "index.html"))
	// 	_ = tmpl.ExecuteTemplate(w, "film-list-element", Film{Title: title, Director: director})
	// }
}
