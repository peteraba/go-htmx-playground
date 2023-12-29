package server

import (
	"crypto/sha512"
	"embed"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	slogfiber "github.com/samber/slog-fiber"

	"github.com/peteraba/go-htmx-playground/pkg/server/middleware"
)

//go:embed assets/*
var assetsFS embed.FS

func Setup(app *fiber.App, logger *slog.Logger, version string) {
	app.Use(recover.New())
	app.Use(middleware.Htmx(version))

	//nolint: exhaustruct
	app.Get("/metrics", monitor.New(monitor.Config{Title: "go|htmx Metrics Page"}))
	app.Use(slogfiber.New(logger))
	app.Use(idempotency.New())
}

func AddStaticHandler(app *fiber.App) {
	const maxAge = 60 * 60

	//nolint: exhaustruct
	app.Use(filesystem.New(filesystem.Config{
		Root: http.FS(assetsFS),
		// PathPrefix:   "/assets",
		Browse:       false,
		NotFoundFile: "404.html",
		MaxAge:       maxAge,
	}))
}

func ListAssets(withIntegrity bool) error {
	list, err := assetsFS.ReadDir("assets")
	if err != nil {
		return fmt.Errorf("failed to read assets directory: %w", err)
	}

	for _, file := range list {
		if withIntegrity {
			content, err := assetsFS.ReadFile(fmt.Sprintf("assets/%s", file.Name()))
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", file.Name(), err)
			}

			sum := sha512.Sum384(content)
			sumSlice := make([]byte, sha512.Size384)

			for i := 0; i < sha512.Size384; i++ {
				sumSlice[i] = sum[i]
			}

			encoded := base64.StdEncoding.EncodeToString(sumSlice)
			fmt.Printf("%s %s\n", file.Name(), encoded) // nolint: forbidigo
		} else {
			fmt.Println(file.Name()) // nolint: forbidigo
		}
	}

	return nil
}
