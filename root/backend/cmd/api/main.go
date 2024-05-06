package main

import (
	"flag"
	"log/slog"
	"os"


	"github.com/ls13g12/hockey-app/root/backend/api"
	"github.com/ls13g12/hockey-app/root/backend/db"
)

var cfg api.Config

func main() {
	flag.StringVar(&cfg.Addr, "addr", ":8080", "HTTP network address")
	flag.StringVar(&cfg.Dsn, "dsn", "mongodb://localhost:27017/hockeydb", "mongodb connection string")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info("attempting to connect to db")
	dbClient, err := db.InitDB(cfg.Dsn)
	if err != nil {
		logger.Error("error connecting to db")
	}

	api.NewApiServer(
		cfg,
		logger,
		dbClient,
	)
}
