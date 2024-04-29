package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/ls13g12/hockey-app/root/backend/api/router"
)


type config struct {
	addr string
}

var cfg config

func NewServer(infoLog *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	router.AddRoutes(mux, infoLog)
	var handler http.Handler = mux
	return handler
}

func main() {
		flag.StringVar(&cfg.addr, "addr", ":8080", "HTTP network address")
		flag.Parse()

		infoLog := slog.New(slog.NewTextHandler(os.Stdout, nil))
		srv := NewServer(infoLog)
	
		httpServer := &http.Server{
			Addr:    cfg.addr,
			Handler: srv,
		}
	
		fmt.Printf("listening on %s\n", httpServer.Addr)
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
}










