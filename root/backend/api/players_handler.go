package api

import (
	"encoding/json"
	"net/http"

	"github.com/ls13g12/hockey-app/root/backend/db"
)

func (a *api) playersGet(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json")

	players, err := db.AllPlayers(a.dbClient)

	if err != nil {
		a.logger.Error("Error %v", err)
	}

	jsonData, err := json.Marshal(players)
	if err != nil {
		a.logger.Error("Error %v", err)
	}
	w.Write(jsonData)
}
