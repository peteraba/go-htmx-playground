package htmx

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	HeaderHxCurrentURL = "Hx-Current-Url"
	HeaderHxRequest    = "Hx-Request"
	HeaderHxTarget     = "Hx-Target"
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

	for k := range headers {
		if k == HeaderHxRequest || k == HeaderHxCurrentURL || k == HeaderHxTarget {
			return true
		}
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
