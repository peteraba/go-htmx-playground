package log

import (
	"log/slog"
)

func Err(e error) slog.Attr {
	return slog.Any("err", e)
}
