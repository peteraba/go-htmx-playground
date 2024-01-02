package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/lib/auth"
	"github.com/peteraba/go-htmx-playground/lib/contenttype"
	"github.com/peteraba/go-htmx-playground/pkg/server/view"
)

func Htmx(buildVersion string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()

		if !contenttype.IsHTML(headers) {
			return c.Next()
		}

		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

		err := c.Next()
		if err != nil {
			return err
		}

		if contenttype.IsHTMX(headers) {
			return nil
		}

		content := c.Response().Body()
		c.Response().ResetBody()

		isAuthenticated := auth.IsAuthenticated(c.Context())
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
