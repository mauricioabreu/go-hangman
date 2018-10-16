package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

var router *mux.Router

func init() {
	configureStoreType("memory", "", "", "")

	router = Router("testWord.txt")
}

func json2GameInfo(bytes []byte, t *testing.T) gameInfoJSON {
	var gameInfo gameInfoJSON
	err := json.Unmarshal(bytes, &gameInfo)
	if err != nil {
		t.Error(err)
	}
	return gameInfo
}

func TestGetInvalidGameID(t *testing.T) {
	r, _ := http.NewRequest("GET", "/games/some-bad-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("Should not be able to request state of a game that doesn't exist. Status Code %d", w.Code)
	}
}

func TestGetInvalidGameIDGuess(t *testing.T) {
	r, _ := http.NewRequest("PUT", "/games/some-bad-id/guesses", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("Should not be able to guess against a game that doesn't exist. Status Code %d", w.Code)
	}
}

func TestGetInvalidGameIDDelete(t *testing.T) {
	r, _ := http.NewRequest("DELETE", "/games/some-bad-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if w.Code != 404 {
		t.Errorf("Should not be able to delete a game that doesn't exist. Status Code %d", w.Code)
	}
}

func TestCreateAndDeleteGame(t *testing.T) {
	// Create game.
	r, _ := http.NewRequest("POST", "/games", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != 204 {
		t.Errorf("Failed to create new game. Status code %d", w.Code)
	}

	// Delete game.
	gameID := strings.Split(w.HeaderMap["Location"][0], "/")[2]
	r, _ = http.NewRequest("DELETE", "/games/"+gameID, nil)
	w = httptest.NewRecorder()
	if w.Code != 200 {
		t.Errorf("Failed to delete game. Status Code %d", w.Code)
	}
}

func TestCreateAndGetGameState(t *testing.T) {
	// Create the game.
	r, _ := http.NewRequest("POST", "/games", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != 204 {
		t.Errorf("Failed to create new game. Status Code %d", w.Code)
	}

	// Get game ID.
	gameID := strings.Split(w.HeaderMap["Location"][0], "/")[2]

	// Get the game Info.
	r, _ = http.NewRequest("GET", "/games/"+gameID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)

	// Check game state
	actualGameInfo := json2GameInfo(w.Body.Bytes(), t)
	expectedGameInfo := gameInfoJSON{gameID, "_____", 3, []string{}, 3}
	if !expectedGameInfo.equals(actualGameInfo) {
		t.Errorf("Expected GameInfo %v, got %v", expectedGameInfo, actualGameInfo)
	}
}

func TestCreateAndPlayGame(t *testing.T) {
	// Create the game.
	r, _ := http.NewRequest("POST", "/games", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != 204 {
		t.Error("Failed to create new game.", w.Code)
	}

	// Get game ID.
	gameID := strings.Split(w.HeaderMap["Location"][0], "/")[2]

	// Guess a correct letter
	r, _ = http.NewRequest("PUT", "/games/"+gameID+"/guesses", nil)
	content := "{\"Guess\":\"t\"}"
	r.Body = ioutil.NopCloser(strings.NewReader(content))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != 201 {
		t.Errorf("Failed to guess a letter. Status Code %d", w.Code)
	}

	// Check game state
	actualGameInfo := json2GameInfo(w.Body.Bytes(), t)
	expectedGameInfo := gameInfoJSON{gameID, "t__t_", 3, []string{"t"}, 3}
	if !expectedGameInfo.equals(actualGameInfo) {
		t.Errorf("Expected GameInfo %v, got %v", expectedGameInfo, actualGameInfo)
	}
}

func TestCreateAndPlayGameWrongGuess(t *testing.T) {
	// Create the game.
	r, _ := http.NewRequest("POST", "/games", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != 204 {
		t.Error("Failed to create new game.", w.Code)
	}

	// Get game ID.
	gameID := strings.Split(w.HeaderMap["Location"][0], "/")[2]

	// Guess a wrong letter
	r, _ = http.NewRequest("PUT", "/games/"+gameID+"/guesses", nil)
	content := "{\"Guess\":\"z\"}"
	r.Body = ioutil.NopCloser(strings.NewReader(content))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != 201 {
		t.Errorf("Failed to guess a letter game %v", w)
	}

	// Check game state
	actualGameInfo := json2GameInfo(w.Body.Bytes(), t)
	expectedGameInfo := gameInfoJSON{gameID, "_____", 2, []string{"z"}, 3}
	if !expectedGameInfo.equals(actualGameInfo) {
		t.Errorf("Expected GameInfo %v, got %v", expectedGameInfo, actualGameInfo)
	}
}

func TestFullGame(t *testing.T) {
	// Create the game.
	r, _ := http.NewRequest("POST", "/games", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != 204 {
		t.Error("Failed to create new game.", w.Code)
	}

	// Get game ID.
	gameID := strings.Split(w.HeaderMap["Location"][0], "/")[2]

	// Guess a correct letter
	r, _ = http.NewRequest("PUT", "/games/"+gameID+"/guesses", nil)
	content := "{\"Guess\":\"t\"}"
	r.Body = ioutil.NopCloser(strings.NewReader(content))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != 201 {
		t.Errorf("Failed to guess a letter game %v", w)
	}
	actualGameInfo := json2GameInfo(w.Body.Bytes(), t)
	expectedGameInfo := gameInfoJSON{gameID, "t__t_", 3, []string{"t"}, 3}
	if !expectedGameInfo.equals(actualGameInfo) {
		t.Errorf("Expected GameInfo %v, got %v", expectedGameInfo, actualGameInfo)
	}

	// Guess another correct letter
	r, _ = http.NewRequest("PUT", "/games/"+gameID+"/guesses", nil)
	content = "{\"Guess\":\"e\"}"
	r.Body = ioutil.NopCloser(strings.NewReader(content))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != 201 {
		t.Errorf("Failed to guess a letter game %v", w)
	}
	actualGameInfo = json2GameInfo(w.Body.Bytes(), t)
	expectedGameInfo = gameInfoJSON{gameID, "te_t_", 3, []string{"t", "e"}, 3}
	if !expectedGameInfo.equals(actualGameInfo) {
		t.Errorf("Expected GameInfo %v, got %v", expectedGameInfo, actualGameInfo)
	}

	// Guess another correct letter
	r, _ = http.NewRequest("PUT", "/games/"+gameID+"/guesses", nil)
	content = "{\"Guess\":\"s\"}"
	r.Body = ioutil.NopCloser(strings.NewReader(content))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, r)
	if w.Code != 201 {
		t.Errorf("Failed to guess a letter game %v", w)
	}
	actualGameInfo = json2GameInfo(w.Body.Bytes(), t)
	expectedGameInfo = gameInfoJSON{gameID, "tests", 3, []string{"t", "e", "s"}, 3}
	if !expectedGameInfo.equals(actualGameInfo) {
		t.Errorf("Expected GameInfo %v, got %v", expectedGameInfo, actualGameInfo)
	}
}
