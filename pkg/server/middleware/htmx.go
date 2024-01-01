package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/lib/auth"
	"github.com/peteraba/go-htmx-playground/lib/htmx"
	"github.com/peteraba/go-htmx-playground/pkg/server/view"
)

func Htmx(buildVersion string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()

		if !htmx.AcceptHTML(headers) {
			return c.Next() // nolint: wrapcheck
		}

		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

		err := c.Next()
		if err != nil {
			return err // nolint: wrapcheck
		}

		if htmx.IsHx(headers) {
			return nil
		}

		content := c.Response().Body()
		c.Response().ResetBody()

		// wrap
		isAuthenticated := c.Context().UserValue(auth.Authenticated).(bool)
		topNav := view.Nav(isAuthenticated)

		version := buildVersion
		if buildVersion == "development" {
			version = time.Now().Format(time.RFC3339Nano)
		}

		component := view.Layout(c.BaseURL(), string(content), topNav, version)
		err = component.Render(c.Context(), c.Response().BodyWriter())

		return err
	}
}
