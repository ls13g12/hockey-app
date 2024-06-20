package server

import (
	"encoding/json"
	"net/http"

	"github.com/ls13g12/hockey-app/src/pkg/db"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ERR_MISSING_FIRST_AND_LAST_NAME = "First and last name required"
)

type PlayerStore interface {
	AllPlayers() ([]db.Player, error)
	GetPlayer(playerID string) (db.Player, error)
	CreatePlayer(player db.Player) error
	UpdatePlayer(player db.Player) error
	DeletePlayer(playerID string) error
}

type PlayerModel struct {
	db *mongo.Database
}

func (pm PlayerModel) AllPlayers() ([]db.Player, error) {
	return db.AllPlayers(pm.db)
}

func (pm PlayerModel) GetPlayer(playerID string) (db.Player, error) {
	return db.GetPlayer(pm.db, playerID)
}

func (pm PlayerModel) CreatePlayer(player db.Player) error {
	return db.CreatePlayer(pm.db, player)
}

func (pm PlayerModel) UpdatePlayer(player db.Player) error {
	return db.UpdatePlayer(pm.db, player)
}

func (pm PlayerModel) DeletePlayer(playerID string) error {
	return db.DeletePlayer(pm.db, playerID)
}

func (s *server) playerGetAll(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json")

	players, err := s.playerStore.AllPlayers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(players)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

func (s *server) playerGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var err error
	w.Header().Set("Content-Type", "application/json")

	player, err := s.playerStore.GetPlayer(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

func (s *server) playerCreate(w http.ResponseWriter, r *http.Request) {
	var player db.Player
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	player.PlayerID = uuid.NewString()

	if player.FirstName == "" || player.LastName == "" {
		http.Error(w, ERR_MISSING_FIRST_AND_LAST_NAME, http.StatusBadRequest)
		return
	}

	if err := s.playerStore.CreatePlayer(player); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *server) playerPut(w http.ResponseWriter, r *http.Request) {
	var player db.Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.playerStore.UpdatePlayer(player); err != nil {
		s.logger.Error("Error %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) playerDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := s.playerStore.DeletePlayer(id); err != nil {
		s.logger.Error("Error %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
