package router

import (
	"log/slog"
	"net/http"

	"github.com/ls13g12/hockey-app/root/backend/api/middleware"
	"github.com/ls13g12/hockey-app/root/backend/api/handlers/healthcheck"
)


func AddRoutes(
	mux *http.ServeMux,
	logger *slog.Logger,	
) {
	middleware := newMiddleware(logger)

	mux.Handle("GET /healthcheck", middleware(healthcheck.Get))
	mux.Handle("/", middleware(notFoundHandler))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func newMiddleware(
	logger *slog.Logger,
) func(h http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.CORSMiddleware(
			middleware.HttpLogger(
				logger, 
				h,
			),
		)
	}
}


