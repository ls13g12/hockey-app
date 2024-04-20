package routers

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/ls13g12/hockey-app/root/backend/pkg/api/middlewares"

	"github.com/ls13g12/hockey-app/root/backend/pkg/api/handlers/healthcheck"
)

func NewServer() http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux)
	var handler http.Handler = mux
	return handler
}

func addRoutes(mux *http.ServeMux) {
	mux.Handle("/healthcheck", middlewares.CORSMiddleware(healthcheck.Index))
	mux.Handle("/", http.NotFoundHandler())
}

func InitApiServer() {
	srv := NewServer()

	httpServer := &http.Server{
		Addr:    net.JoinHostPort("localhost", "8080"),
		Handler: srv,
	}

	fmt.Printf("listening on %s\n", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
	}
	
}
