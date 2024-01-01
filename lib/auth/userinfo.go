package auth

import (
	"github.com/valyala/fasthttp"
)

const (
	Authenticated = "authenticated"
	SubjectKey    = "Subject"
	Name          = "Name"
)

func IsAuthenticated(ctx *fasthttp.RequestCtx) bool {
	result, ok := ctx.UserValue(Authenticated).(bool)
	if !ok {
		return false
	}

	return result
}

func GetString(ctx *fasthttp.RequestCtx, key string) string {
	result, ok := ctx.UserValue(key).(string)
	if !ok {
		return ""
	}

	return result
}
