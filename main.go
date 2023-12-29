package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v2"

	"github.com/peteraba/go-htmx-playground/pkg/colors"
	"github.com/peteraba/go-htmx-playground/pkg/films"
	"github.com/peteraba/go-htmx-playground/pkg/home"
	"github.com/peteraba/go-htmx-playground/pkg/notifications"
	notificationsService "github.com/peteraba/go-htmx-playground/pkg/notifications/service"
	"github.com/peteraba/go-htmx-playground/pkg/server"
)

var Version = "development" // nolint:gochecknoglobals

const maxListLength = 10

// nolint: exhaustruct, gomnd, varnamelen, wrapcheck, funlen
func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name: "assets",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "integrity",
						Aliases: []string{"i"},
						Usage:   "calculate the integrity of the assets",
					},
				},
				Action: func(c *cli.Context) error {
					return server.ListAssets(c.Bool("integrity"))
				},
			},
			{
				Name: "version",
				Action: func(c *cli.Context) error {
					fmt.Printf("Current version: %s\n", Version) // nolint: forbidigo

					return nil
				},
			},
			{
				Name:  "server",
				Usage: "Start the server",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   8000,
					},
				},
				Action: func(c *cli.Context) error {
					notifier := notificationsService.NewNotifier(logger)
					// nolint: exhaustruct
					app := fiber.New(fiber.Config{
						Immutable: true,
					})

					server.Setup(app, logger, Version)
					notifications.Setup(app, logger, notifier)
					home.Setup(app)
					colors.Setup(app)
					films.Setup(app, logger, notifier, maxListLength, Version)
					server.AddStaticHandler(app)

					err := app.Listen(fmt.Sprintf(":%d", c.Int("port")))
					if err != nil {
						logger.Error(err.Error())
						os.Exit(1)
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Error(err.Error())
	}
}
