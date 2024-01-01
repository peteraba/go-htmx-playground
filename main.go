package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v2"

	"github.com/peteraba/go-htmx-playground/lib/auth"
	"github.com/peteraba/go-htmx-playground/pkg/colors"
	"github.com/peteraba/go-htmx-playground/pkg/dashboard"
	"github.com/peteraba/go-htmx-playground/pkg/home"
	"github.com/peteraba/go-htmx-playground/pkg/movies"
	"github.com/peteraba/go-htmx-playground/pkg/notifications"
	notificationsService "github.com/peteraba/go-htmx-playground/pkg/notifications/service"
	"github.com/peteraba/go-htmx-playground/pkg/server"
)

var Version = "development" // nolint:gochecknoglobals

const maxListLength = 10

// nolint: exhaustruct, gomnd, varnamelen, wrapcheck, funlen
func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

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
					&cli.StringFlag{
						Name:    "zitadel",
						Aliases: []string{"z"},
						Value:   "localhost",
						EnvVars: []string{"ZITADEL"},
					},
					&cli.StringFlag{
						Name:    "zitadel-key",
						Aliases: []string{"zk"},
						Value:   "",
						EnvVars: []string{"ZITADEL_MASTERKEY"},
					},
					&cli.StringFlag{
						Name:    "zitadel-client-id",
						Aliases: []string{"zci"},
						Value:   "",
						EnvVars: []string{"ZITADEL_CLIENT_ID"},
					},
					&cli.StringFlag{
						Name:    "zitadel-port",
						Aliases: []string{"zp"},
						Value:   "443",
						EnvVars: []string{"ZITADEL_PORT"},
					},
					&cli.BoolFlag{
						Name:    "zitadel-insecure",
						Aliases: []string{"zi"},
						Value:   false,
						EnvVars: []string{"ZITADEL_INSECURE"},
					},
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   8000,
					},
				},
				Action: func(c *cli.Context) error {
					zitadelOptions := auth.Options{
						Domain:      c.String("zitadel"),
						ClientID:    c.String("zitadel-client-id"),
						Key:         c.String("zitadel-key"),
						Port:        c.String("zitadel-port"),
						Insecure:    c.Bool("zitadel-insecure"),
						RedirectURI: "http://localhost:8000/auth/callback",
					}

					notifier := notificationsService.NewNotifier(logger)
					// nolint: exhaustruct
					app := fiber.New(fiber.Config{
						Immutable: true,
					})

					requireAuthHandler := auth.Setup(app, zitadelOptions, logger)
					server.Setup(app, logger, Version)
					notifications.Setup(app, logger, notifier)
					home.Setup(app, logger)
					colors.Setup(app)
					movies.Setup(app, logger, notifier, maxListLength, Version)

					// protected routes
					dashboard.Setup(app, requireAuthHandler, logger)

					// static server
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
