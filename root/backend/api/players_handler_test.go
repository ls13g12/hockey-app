package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ls13g12/hockey-app/root/backend/db"
	"github.com/stretchr/testify/assert"
)

type MockedPlayerModel struct {
	players []db.Player
	err error
}

func (mpm MockedPlayerModel) AllPlayers() ([]db.Player, error) {
	return mpm.players, mpm.err
}

func TestPlayersGet(t *testing.T) {

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

	testApi.playersGet(rr, r)

	rs := rr.Result()

	assert.Equal(t, rs.StatusCode, http.StatusOK)
	
	var players []db.Player
	err = json.Unmarshal(rr.Body.Bytes(), &players)

	assert.Equal(t, len(players), 2)
}

