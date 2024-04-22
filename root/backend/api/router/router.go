package router

import (
	"net/http"

	"github.com/ls13g12/hockey-app/root/backend/api/middleware"

	"github.com/ls13g12/hockey-app/root/backend/api/handlers/healthcheck"
)


func AddRoutes(mux *http.ServeMux) {
	mux.Handle("/healthcheck", middleware.CORSMiddleware(healthcheck.Get))
	mux.Handle("/", http.NotFoundHandler())
}


