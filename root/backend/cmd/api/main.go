package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/alexedwards/scs/mongodbstore"
	"github.com/alexedwards/scs/v2"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ls13g12/hockey-app/root/backend/api"
	"github.com/ls13g12/hockey-app/root/backend/db"
)

var cfg api.Config

func main() {
	flag.StringVar(&cfg.Addr, "addr", ":8080", "HTTP network address")
	flag.StringVar(&cfg.Dsn, "dsn", "mongodb://localhost:27017/test", "mongodb connection string")
	flag.StringVar(&cfg.Mode, "mode", "dev", "app mode - dev or prod")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger.Info(fmt.Sprintf("running app in %s mode", cfg.Mode))
	
	logger.Info("attempting to connect to db")
	dbClient, err := db.InitDB(cfg.Dsn)
	if err != nil {
		logger.Error("error connecting to db")
	}

	var db *mongo.Database 
	if cfg.Mode == "prod" {
		db = dbClient.Database("hockeydb")
	} else {
		db = dbClient.Database("test")
	}

	sessionManager := scs.New()
	sessionManager.Store = mongodbstore.New(db)
	sessionManager.Lifetime = 1 * time.Hour
	// sessionManager.Cookie.Persist = true
	// sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	// sessionManager.Cookie.Secure = false

	api.NewApiServer(
		cfg,
		logger,
		db,
		sessionManager,
	)
}
