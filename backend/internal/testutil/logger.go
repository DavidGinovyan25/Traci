package testutil

import (
	"io"
	"log/slog"
)

func Logger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}
