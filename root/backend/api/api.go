package api

import (
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/justinas/alice"
	"github.com/ls13g12/hockey-app/root/backend/token"
	"github.com/ls13g12/hockey-app/root/backend/util"
	"go.mongodb.org/mongo-driver/mongo"
)



type api struct {
	logger 					*slog.Logger
	tokenMaker 			token.Maker
	playerStore 		PlayerStore
	userStore				UserStore
	sessionManager 	*scs.SessionManager
}

func NewApiServer(cfg util.Config, logger *slog.Logger, tokenMaker token.Maker, db *mongo.Database, sm *scs.SessionManager) {
	a := &api{
		logger: logger,
		tokenMaker: tokenMaker,
		playerStore: PlayerModel{db: db},
		userStore: UserModel{db: db},
		sessionManager: sm,
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

	dynamic := alice.New(a.sessionManager.LoadAndSave)
	protected := dynamic.Append(a.isAuthenticated)

	mux.Handle("POST /user/login", dynamic.ThenFunc(a.userLogin))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(a.userSignup))

	mux.Handle("GET /players", dynamic.ThenFunc(a.playerGetAll))
	mux.Handle("GET /players/{id}", dynamic.ThenFunc(a.playerGet))
	mux.Handle("CREATE /players", protected.ThenFunc(a.playerCreate))
	mux.Handle("PUT /players", protected.ThenFunc(a.playerPut))
	mux.Handle("DELETE /players/{id}", protected.ThenFunc(a.playerDelete))

	mux.HandleFunc("/", notFoundHandler)

	standard := alice.New(a.corsHeaders, a.requestLogger)
	
	return standard.Then(mux)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}


