package server

import (
	"log/slog"
	"net/http"

	"github.com/ls13g12/hockey-app/src/pkg/common"
	"github.com/ls13g12/hockey-app/src/pkg/token"

	"github.com/alexedwards/scs/v2"
	"github.com/justinas/alice"
	"go.mongodb.org/mongo-driver/mongo"
)

type server struct {
	logger         *slog.Logger
	tokenMaker     token.Maker
	playerStore    PlayerStore
	userStore      UserStore
	sessionManager *scs.SessionManager
}

func NewServer(cfg common.ServerConfig, logger *slog.Logger, tokenMaker token.Maker, db *mongo.Database, sm *scs.SessionManager) {
	s := &server{
		logger:         logger,
		tokenMaker:     tokenMaker,
		playerStore:    PlayerModel{db: db},
		userStore:      UserModel{db: db},
		sessionManager: sm,
	}

	httpServer := &http.Server{
		Addr:    cfg.Addr,
		Handler: s.addRoutes(),
	}

	s.logger.Info("api server started", slog.String("Addr", cfg.Addr))
	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.logger.Error("error listening and serving", slog.Any("Error", err))
	}
}

func (s *server) addRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthcheck", s.healthcheckGet)

	dynamic := alice.New(s.sessionManager.LoadAndSave)
	protected := dynamic.Append(s.isAuthenticated)

	mux.Handle("POST /user/login", dynamic.ThenFunc(s.userLogin))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(s.userSignup))

	mux.Handle("GET /players", dynamic.ThenFunc(s.playerGetAll))
	mux.Handle("GET /players/{id}", dynamic.ThenFunc(s.playerGet))
	mux.Handle("CREATE /players", protected.ThenFunc(s.playerCreate))
	mux.Handle("PUT /players", protected.ThenFunc(s.playerPut))
	mux.Handle("DELETE /players/{id}", protected.ThenFunc(s.playerDelete))

	mux.HandleFunc("/", notFoundHandler)

	standard := alice.New(s.corsHeaders, s.requestLogger)

	return standard.Then(mux)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
