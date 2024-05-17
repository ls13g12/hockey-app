package api

import (
	"encoding/json"
	"net/http"

	"github.com/ls13g12/hockey-app/root/backend/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlayerStore interface {
	AllPlayers() ([]db.Player, error)
}

type PlayerModel struct {
	db *mongo.Database
}

func (pm PlayerModel) AllPlayers() ([]db.Player, error) {
	return db.AllPlayers(pm.db)
}

func (a *api) playersGet(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json")

	players, err := a.playerStore.AllPlayers()

	if err != nil {
		a.logger.Error("Error %v", err)
	}

	jsonData, err := json.Marshal(players)
	if err != nil {
		a.logger.Error("Error %v", err)
	}
	w.Write(jsonData)
}
