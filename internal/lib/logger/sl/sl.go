package sl

import "log/slog"

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func Info(msg string) slog.Attr {
	return slog.Attr{
		Key:   "info",
		Value: slog.StringValue(msg),
	}
}
