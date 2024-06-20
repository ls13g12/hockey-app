package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/alexedwards/scs/mongodbstore"
	"github.com/alexedwards/scs/v2"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ls13g12/hockey-app/src/pkg/server"
	"github.com/ls13g12/hockey-app/src/pkg/common"
	"github.com/ls13g12/hockey-app/src/pkg/db"
	"github.com/ls13g12/hockey-app/src/pkg/token"
)

var cfg common.ServerConfig

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	flag.StringVar(&cfg.Addr, "addr", ":8080", "HTTP network address")
	flag.StringVar(&cfg.Dsn, "dsn", "mongodb://localhost:27017/test", "mongodb connection string")
	flag.StringVar(&cfg.Mode, "mode", "dev", "app mode - dev or prod")
	privateKeyHex := flag.String("private-key", "2b0ac9e5d6d77aff5b55509f5d3dbca1249cb088ebaee71843e974b2e889661ec729ed1319b8074292124748ac08e867d04f4131dcef740df49ade1c5c437e52", "Ed25519 private key in hex format")
	flag.Parse()

	logger.Info(fmt.Sprintf("running app in %s mode", cfg.Mode))

	var err error
	cfg.TokenPrivateKey, err = hex.DecodeString(*privateKeyHex)
	if err != nil {
		logger.Error("invalid private key format - exiting", slog.Any("Error", err))
		return
	}
	if len(cfg.TokenPrivateKey) != ed25519.PrivateKeySize {
		logger.Error("invalid private key size ", slog.Int("Expected bytes", ed25519.PrivateKeySize))
		return
	}

	cfg.TokenPublicKey = ed25519.PublicKey(cfg.TokenPrivateKey[32:])

	tokenMaker, err := token.NewPasetoMaker(cfg.TokenPrivateKey, cfg.TokenPublicKey)
	if err != nil {
		logger.Error("cannot create token maker", slog.Any("Error", err))
		return
	}

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

	sessionManager := scs.New()
	sessionManager.Store = mongodbstore.New(db)
	sessionManager.Lifetime = 1 * time.Hour
	// sessionManager.Cookie.Persist = true
	// sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	// sessionManager.Cookie.Secure = false

	server.NewServer(
		cfg,
		logger,
		tokenMaker,
		db,
		sessionManager,
	)
}
