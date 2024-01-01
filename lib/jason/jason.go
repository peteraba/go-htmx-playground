package jason

import (
	"reflect"

	"github.com/gofiber/fiber/v2"

	"github.com/peteraba/go-htmx-playground/lib/pagination"
)

type Response struct {
	Query map[string]interface{} `json:"query,omitempty"`
	Self  string                 `json:"self"`
	First string                 `json:"first,omitempty"`
	Prev  string                 `json:"prev,omitempty"`
	Next  string                 `json:"next,omitempty"`
	Last  string                 `json:"last,omitempty"`
	Items interface{}            `json:"items,omitempty"`
}

func SendList(c *fiber.Ctx, items interface{}, p pagination.Pagination) error {
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	val := reflect.ValueOf(items)
	if val.Kind() != reflect.Slice || val.Len() == 0 {
		items = nil
	}

	return c.JSON(Response{
		Query: p.Query,
		Self:  p.SelfLink(),
		First: p.FirstLink(),
		Prev:  p.PrevLink(),
		Next:  p.NextLink(),
		Last:  p.LastLink(),
		Items: items,
	})
}
