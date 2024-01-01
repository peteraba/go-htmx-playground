package auth

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication"
	openid "github.com/zitadel/zitadel-go/v3/pkg/authentication/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"

	"github.com/peteraba/go-htmx-playground/lib/log"
)

type Options struct {
	Domain      string
	ClientID    string
	Key         string
	Port        string
	Insecure    bool
	RedirectURI string
}

func Setup(app *fiber.App, options Options, logger *slog.Logger) fiber.Handler {
	logger.Debug("Zitadel connection",
		slog.String("domain", options.Domain),
		slog.String("clientID", options.ClientID),
		slog.String("key", options.Key),
		slog.String("port", options.Port),
		slog.Bool("insecure", options.Insecure),
		slog.String("redirectURL", options.RedirectURI),
	)

	if options.Domain == "" || options.ClientID == "" || options.Key == "" || options.RedirectURI == "" {
		logger.Error("Zitadel connection is not configured")
		os.Exit(1)
	}

	var zitadelClient *zitadel.Zitadel
	if options.Insecure && options.Port != "" {
		zitadelClient = zitadel.New(options.Domain, zitadel.WithInsecure(options.Port))
	} else {
		zitadelClient = zitadel.New(options.Domain)
	}

	ctx := context.Background()
	initializer := openid.DefaultAuthentication(options.ClientID, options.RedirectURI, options.Key)
	authenticator, err := authentication.New(ctx, zitadelClient, options.Key, initializer)
	if err != nil {
		logger.Error("zitadel sdk could not initialize", log.Err(err))
		os.Exit(1)
	}

	interceptor := authentication.Middleware(authenticator)
	app.Use(adaptor.HTTPMiddleware(interceptor.CheckAuthentication()))

	app.Use(func(c *fiber.Ctx) error {
		authCtx := interceptor.Context(c.Context())
		if authCtx == nil {
			c.Context().SetUserValue(Authenticated, false)

			return c.Next()
		}

		userInfo := authCtx.UserInfo
		c.Context().SetUserValue(SubjectKey, userInfo.Subject)
		c.Context().SetUserValue(Name, userInfo.Name)
		c.Context().SetUserValue(Authenticated, true)

		if err != nil {
			return fmt.Errorf("could not bind user info: %w", err)
		}

		return c.Next()
	})

	app.All("/auth/*", adaptor.HTTPHandler(authenticator))

	return adaptor.HTTPMiddleware(interceptor.RequireAuthentication())
}
