package main

import (
	"flag"
	"log/slog"
	"os"
	"github.com/ls13g12/hockey-app/root/backend/api"
)

type application struct {
	logger *slog.Logger
}

var cfg api.Config

func main() {
	flag.StringVar(&cfg.Addr, "addr", ":8080", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	api.NewApiServer(
		cfg,
		logger,
	)
}
