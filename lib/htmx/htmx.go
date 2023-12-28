package htmx

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/pkg/app/view"
)

const (
	HeaderHxRequest = "Hx-Request"
	HeaderHxTarget  = "Hx-Target"
)

func IsHx(headers map[string][]string) bool {
	if len(headers) == 0 {
		return false
	}

	v, ok := headers[HeaderHxRequest]
	if ok && len(v) > 0 {
		return v[0] == "true"
	}

	return false
}

func GetTarget(headers map[string][]string) string {
	v, ok := headers[HeaderHxTarget]
	if ok && len(v) > 0 {
		return v[0]
	}

	return ""
}

func AcceptHTML(headers map[string][]string) bool {
	if len(headers) == 0 {
		return false
	}

	acceptHeaders, ok := headers[fiber.HeaderAccept]
	if !ok {
		return false
	}

	for _, acceptHeader := range acceptHeaders {
		for _, elem := range strings.Split(acceptHeader, ",") {
			if elem == fiber.MIMETextHTML {
				return true
			}
		}
	}

	return false
}

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()
		if !AcceptHTML(headers) {
			return c.Next() // nolint: wrapcheck
		}

		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)

		err := c.Next()
		if err != nil {
			return err // nolint: wrapcheck
		}

		if IsHx(headers) {
			return nil
		}

		content := c.Response().Body()
		c.Response().ResetBody()

		// wrap
		topNav := view.NewTopNav(c.Path())
		navComponent := topNav.Nav()

		component := view.Layout(c.BaseURL(), string(content), navComponent)
		err = component.Render(c.Context(), c.Response().BodyWriter())

		return err
	}
}
