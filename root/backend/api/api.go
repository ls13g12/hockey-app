package api

import (
	"log/slog"
	"net/http"

	"github.com/justinas/alice"
)

type Config struct {
	Addr string
}

type api struct {
	logger *slog.Logger
}

func NewApiServer(cfg Config, logger *slog.Logger) {
	a := &api{
		logger: logger,
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

	mux.HandleFunc("GET /healthcheck", healthcheckGet)
	mux.HandleFunc("/", notFoundHandler)

	standard := alice.New(a.corsHeaders, a.requestLogger)
	
	return standard.Then(mux)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}


