package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/ls13g12/hockey-app/root/backend/api/middleware"
	"github.com/ls13g12/hockey-app/root/backend/api/router"
)


func NewServer() http.Handler {
	mux := http.NewServeMux()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	router.AddRoutes(mux)
	var handler http.Handler = middleware.HttpLogger(mux, logger)
	return handler
}

func main() {
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










