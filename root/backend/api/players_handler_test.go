package api

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/ls13g12/hockey-app/root/backend/db"
	"github.com/stretchr/testify/assert"
)

type MockedPlayerModel struct {
	players []db.Player
	player db.Player
	err error
}

func (mpm MockedPlayerModel) AllPlayers() ([]db.Player, error) {
	return mpm.players, mpm.err
}

func (mpm MockedPlayerModel) GetPlayer(playerID string) (db.Player, error) {
	return mpm.player, mpm.err
}

func (mpm MockedPlayerModel) CreatePlayer(player db.Player) error {
	return mpm.err
}

func (mpm MockedPlayerModel) UpdatePlayer(player db.Player) error {
	return mpm.err
}

func (mpm MockedPlayerModel) DeletePlayer(playerID string) error {
	return mpm.err
}

func TestPlayersGetAll(t *testing.T) {

	mockedPlayers := []db.Player{
		{
			PlayerID: "1",
		},
		{
			PlayerID: "2",
		},
	}

	testApi := api{
		logger: nil,
		playerStore: MockedPlayerModel{players: mockedPlayers},
	}

	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/players", nil)
	if err != nil {
			t.Fatal(err)
	}

	testApi.playerGetAll(rr, r)

	rs := rr.Result()

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	
	var players []db.Player
	err = json.Unmarshal(rr.Body.Bytes(), &players)

	assert.Equal(t, len(players), 2)
}

func TestPlayerCreate(t *testing.T) {
	newPlayer := db.Player{
			PlayerID: "3",
			FirstName: "John",
			LastName: "Doe",
			Nickname: "Johnny",
			HomeShirtNumber: 10,
			DateOfBirth: time.Now().AddDate(-25, 0, 0),
	}

	testApi := api{
			logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
			playerStore: MockedPlayerModel{err: nil},
	}

	rr := httptest.NewRecorder()
	body, err := json.Marshal(newPlayer)
	if err != nil {
			t.Fatal(err)
	}

	r, err := http.NewRequest(http.MethodPost, "/players", bytes.NewBuffer(body))
	if err != nil {
			t.Fatal(err)
	}

	testApi.playerCreate(rr, r)

	rs := rr.Result()
	assert.Equal(t, rs.StatusCode, http.StatusCreated)
}

func TestPlayerUpdate(t *testing.T) {
	updatedPlayer := db.Player{
			PlayerID: "1",
			FirstName: "Jane",
			LastName: "Doe",
			Nickname: "Janie",
			HomeShirtNumber: 15,
			DateOfBirth: time.Now().AddDate(-22, 0, 0),
	}

	testApi := api{
			logger: nil,
			playerStore: MockedPlayerModel{err: nil},
	}

	rr := httptest.NewRecorder()
	body, err := json.Marshal(updatedPlayer)
	if err != nil {
			t.Fatal(err)
	}

	r, err := http.NewRequest(http.MethodPut, "/players", bytes.NewBuffer(body))
	if err != nil {
			t.Fatal(err)
	}

	testApi.playerPut(rr, r)

	rs := rr.Result()
	assert.Equal(t, rs.StatusCode, http.StatusOK)
}

func TestPlayerDelete(t *testing.T) {
	testApi := api{
			logger: nil,
			playerStore: MockedPlayerModel{err: nil},
	}

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodDelete, "/players/1", nil)
	if err != nil {
			t.Fatal(err)
	}

	testApi.playerDelete(rr, r)

	rs := rr.Result()
	assert.Equal(t, rs.StatusCode, http.StatusOK)
}
