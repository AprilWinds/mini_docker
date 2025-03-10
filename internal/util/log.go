package util

import (
	"log/slog"
	"os"
)

func LogAndExit(msg string, err error) {
	if err != nil {
		slog.Error(msg, "err", err)
	} else {
		slog.Error(msg)
	}
	os.Exit(1)
}
