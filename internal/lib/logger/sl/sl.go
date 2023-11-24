package sl

import (
	"log/slog"
)

//для добавления ошибки в log, чтобы коротко писать slog.Err(err)
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}