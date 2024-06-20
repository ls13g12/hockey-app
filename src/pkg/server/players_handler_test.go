package server

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ls13g12/hockey-app/src/pkg/common"
	"github.com/ls13g12/hockey-app/src/pkg/db"

	"github.com/stretchr/testify/assert"
)

type MockedPlayerModel struct {
	players []db.Player
	player  db.Player
	err     error
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

func generateMockPlayers(numPlayers int) []db.Player {
	players := make([]db.Player, numPlayers)
	for i := 0; i < numPlayers; i++ {
		players[i] = db.Player{
			PlayerID:        common.GenerateRandomString(10),
			FirstName:       common.GenerateRandomString(5),
			LastName:        common.GenerateRandomString(7),
			Nickname:        common.GenerateRandomString(4),
			HomeShirtNumber: rand.Intn(99) + 1,
			AwayShirtNumber: rand.Intn(99) + 1,
			DateOfBirth:     common.GenerateRandomDate(),
			Created:         time.Now(),
		}
	}
	return players
}

func generateCreatePlayerPayload(
	firstNameRequired bool,
	lastNameRequired bool,
) db.Player {
	var player db.Player

	if firstNameRequired {
		player.FirstName = common.GenerateRandomString(5)
	}

	if lastNameRequired {
		player.LastName = common.GenerateRandomString(7)
	}

	player.Nickname = common.GenerateRandomString(4)
	player.HomeShirtNumber = rand.Intn(99) + 1
	player.AwayShirtNumber = rand.Intn(99) + 1
	player.DateOfBirth = common.GenerateRandomDate()
	player.Created = time.Now()

	return player
}

func TestPlayersGetAll(t *testing.T) {

	testcases := []struct {
		testName       string
		mockedPlayers  []db.Player
		wantStatusCode int
	}{
		{
			testName:       "GetAllPlayers - Success",
			mockedPlayers:  generateMockPlayers(5),
			wantStatusCode: http.StatusOK,
		},
		{
			testName:       "GetAllPlayersEmpty - Success",
			mockedPlayers:  nil,
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testcases {
		testApi := server{
			playerStore: MockedPlayerModel{players: tc.mockedPlayers},
		}

		rr := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodGet, "/players", nil)
		if err != nil {
			t.Fatal(err)
		}

		testApi.playerGetAll(rr, r)

		rs := rr.Result()

		assert.Equal(t, rs.StatusCode, tc.wantStatusCode)

		var players []db.Player
		err = json.Unmarshal(rr.Body.Bytes(), &players)

		assert.Equal(t, len(players), len(tc.mockedPlayers))
	}
}

func TestPlayerCreate(t *testing.T) {
	testcases := []struct {
		testName         string
		payload          db.Player
		wantStatusCode   int
		wantErrorMessage string
	}{
		{
			testName:       "CreateNewPlayer - Success",
			payload:        generateCreatePlayerPayload(true, true),
			wantStatusCode: http.StatusCreated,
		},
		{
			testName:         "CreateNewPlayer - MissingFirstName",
			payload:          generateCreatePlayerPayload(false, true),
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: ERR_MISSING_FIRST_AND_LAST_NAME,
		},
		{
			testName:         "CreateNewPlayer - MissingLastName",
			payload:          generateCreatePlayerPayload(true, false),
			wantStatusCode:   http.StatusBadRequest,
			wantErrorMessage: ERR_MISSING_FIRST_AND_LAST_NAME,
		},
	}

	for _, tc := range testcases {
		testApi := server{
			logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
			playerStore: MockedPlayerModel{err: nil},
		}

		rr := httptest.NewRecorder()
		body, err := json.Marshal(tc.payload)
		if err != nil {
			t.Fatal(err)
		}

		r, err := http.NewRequest(http.MethodPost, "/players", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}

		testApi.playerCreate(rr, r)

		rs := rr.Result()
		assert.Equal(t, rs.StatusCode, tc.wantStatusCode)
		if rs.StatusCode > 299 {
			assert.Equal(t, strings.TrimSpace(rr.Body.String()), tc.wantErrorMessage)
		}
	}
}

func TestPlayerUpdate(t *testing.T) {
	updatedPlayer := db.Player{
		PlayerID:        "1",
		FirstName:       "Jane",
		LastName:        "Doe",
		Nickname:        "Janie",
		HomeShirtNumber: 15,
		DateOfBirth:     time.Now().AddDate(-22, 0, 0),
	}

	testApi := server{
		logger:      nil,
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
	testApi := server{
		logger:      nil,
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
