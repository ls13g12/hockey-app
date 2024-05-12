package api

import (
	"log/slog"
	"net/http"

	"github.com/justinas/alice"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	Addr string
	Dsn  string
}

type api struct {
	logger *slog.Logger
	playerStore PlayerStore
}

func NewApiServer(cfg Config, logger *slog.Logger, dbClient *mongo.Client) {
	a := &api{
		logger: logger,
		playerStore: PlayerModel{dbClient: dbClient},
	}
	
	httpServer := &http.Server{
		Addr:    cfg.Addr,
		Handler: a.addRoutes(),
	}

	a.logger.Info("api server started", slog.String("Addr", cfg.Addr))
	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		a.logger.Error("error listening and serving", slog.Any("Error", err))
	}
}



func (a *api) addRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthcheck", a.healthcheckGet)

	mux.HandleFunc("GET /players", a.playersGet)

	mux.HandleFunc("/", notFoundHandler)

	standard := alice.New(a.corsHeaders, a.requestLogger)
	
	return standard.Then(mux)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}


