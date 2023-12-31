package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/urfave/cli/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication"
	openid "github.com/zitadel/zitadel-go/v3/pkg/authentication/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"

	"github.com/peteraba/go-htmx-playground/lib/log"
	"github.com/peteraba/go-htmx-playground/pkg/colors"
	"github.com/peteraba/go-htmx-playground/pkg/dashboard"
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
					&cli.StringFlag{
						Name:    "zitadel",
						Aliases: []string{"z"},
						Value:   "localhost",
					},
					&cli.StringFlag{
						Name:    "zitadel-client-id",
						Aliases: []string{"zci"},
						Value:   "",
					},
					&cli.StringFlag{
						Name:    "zitadel-key",
						Aliases: []string{"zk"},
						Value:   "",
					},
					&cli.IntFlag{
						Name:    "zitadel-port",
						Aliases: []string{"zp"},
						Value:   8080,
					},
					&cli.BoolFlag{
						Name:    "zitadel-insecure",
						Aliases: []string{"zi"},
						Value:   false,
					},
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   8000,
					},
				},
				Action: func(c *cli.Context) error {
					zDomain := c.String("zitadel")
					zClientID := c.String("zitadel-client-id")
					zKey := c.String("zitadel-key")
					zPort := c.Int("zitadel-port")
					zInsecure := c.Bool("zitadel-insecure")
					// redirectURI := fmt.Sprintf("http://localhost:%d/dashboard", c.Int("port"))
					redirectURI := fmt.Sprintf("http://localhost:%d/auth/callback", c.Int("port"))

					logger.Info("Zitadel connection",
						slog.String("domain", zDomain),
						slog.String("clientID", zClientID),
						slog.String("key", zKey),
						slog.Int("port", zPort),
						slog.Bool("insecure", zInsecure),
						slog.String("redirectURL", redirectURI),
					)

					if zDomain == "" || zClientID == "" || zKey == "" {
						logger.Error("Zitadel connection is not configured")
						os.Exit(1)
					}

					var zitadelClient *zitadel.Zitadel
					if c.Bool("zitadel-insecure") && c.String("zitadel-port") != "" {
						zitadelClient = zitadel.New(zDomain, zitadel.WithInsecure(c.String("zitadel-port")))
					} else {
						zitadelClient = zitadel.New(zDomain)
					}

					authenticator, err := authentication.New(c.Context, zitadelClient, zKey, openid.DefaultAuthentication(zClientID, redirectURI, zKey))
					if err != nil {
						slog.Error("zitadel sdk could not initialize", log.Err(err))
						os.Exit(1)
					}

					interceptor := authentication.Middleware(authenticator)

					notifier := notificationsService.NewNotifier(logger)
					// nolint: exhaustruct
					app := fiber.New(fiber.Config{
						Immutable: true,
					})

					// unprotected routes
					server.Setup(app, logger, Version)
					notifications.Setup(app, logger, notifier)
					home.Setup(app, interceptor, logger)
					colors.Setup(app)
					films.Setup(app, logger, notifier, maxListLength, Version)

					// authentication
					app.All("/auth/*", adaptor.HTTPHandler(authenticator))

					// protected routes
					userInfoRetriever := func(ctx context.Context) ([]byte, error) {
						authCtx := interceptor.Context(ctx)
						logger.Info("Auth Context", slog.Any("authCtx", authCtx))
						logger.Info("Retrieving user info", slog.Any("user info", authCtx.UserInfo))
						data, err := json.MarshalIndent(authCtx.UserInfo, "", "  ")
						if err != nil {
							return nil, fmt.Errorf("error marshalling profile response: %w", err)
						}

						return data, nil
					}
					dashboard.Setup(app, interceptor, userInfoRetriever, logger)

					// static server
					server.AddStaticHandler(app)

					err = app.Listen(fmt.Sprintf(":%d", c.Int("port")))
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
