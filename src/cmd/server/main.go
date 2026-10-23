package main

import (
	"log/slog"
	"os"

	"github.com/calebportodev/golang/hare-st/src/internal/server"
)

var version = "0.0.1"

func init() {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				},
			),
		),
	)
}

func main() {
	slog.Info("Starting server...", "version", version)
	server.Start(":8379")
}
