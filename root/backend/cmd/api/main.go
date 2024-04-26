package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/ls13g12/hockey-app/root/backend/api/router"
)


func NewServer(
	logger *slog.Logger,
) http.Handler {
	mux := http.NewServeMux()

	router.AddRoutes(mux, logger)
	var handler http.Handler = mux
	return handler
}

func main() {
		logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
		srv := NewServer(logger)
	
		httpServer := &http.Server{
			Addr:    net.JoinHostPort("localhost", "8080"),
			Handler: srv,
		}
	
		fmt.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
}










