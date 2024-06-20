package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/ls13g12/hockey-app/src/pkg/common"
	"github.com/ls13g12/hockey-app/src/pkg/db"
	"github.com/ls13g12/hockey-app/src/pkg/tui"
	"go.mongodb.org/mongo-driver/mongo"
)

var cfg common.TuiAppConfig


func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	flag.StringVar(&cfg.Dsn, "dsn", "mongodb://localhost:27017/test", "mongodb connection string")
	flag.StringVar(&cfg.Mode, "mode", "dev", "app mode - dev or prod")
	flag.Parse()

	logger.Info(fmt.Sprintf("running app in %s mode", cfg.Mode))


	logger.Info("attempting to connect to db")
	dbClient, err := db.InitDB(cfg.Dsn)
	if err != nil {
		logger.Error("error connecting to db", slog.Any("Error", err))
	}

	var db *mongo.Database
	if cfg.Mode == "prod" {
		db = dbClient.Database("hockeydb")
	} else {
		db = dbClient.Database("test")
	}


	tui.NewTuiApp(
		cfg,
		logger,
		db,
	)
}



