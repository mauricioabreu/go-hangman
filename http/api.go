package main

import (
	"encoding/json"
	"go-hangman/game"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type gameInfoJSON struct {
	ID             string   `json:"id"`
	TurnsLeft      int      `json:"turns_left"`
	Used           []string `json:"used"`
	AvailableHints int      `json:"available_hints"`
}

type userGuess struct {
	Guess string
}

func newGame(w http.ResponseWriter, r *http.Request) {
	words := []string{
		"apple",
		"banana",
		"orange",
	}
	choosenWord := hangman.PickWord(words)
	game := hangman.NewGame(3, choosenWord)
	hangman.CreateGame(game)
	w.Header().Set("Location", strings.Join([]string{r.Host, "games", game.ID}, "/"))
}

func retrieveGameInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	game, err := hangman.RetrieveGame(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	responseJSON := gameInfoJSON{
		ID:             game.ID,
		TurnsLeft:      game.TurnsLeft,
		Used:           game.Used,
		AvailableHints: game.AvailableHints,
	}
	buff, error := json.MarshalIndent(responseJSON, "", "\t")
	if error != nil {
		log.Fatal("Could not serialize game")
	}

	w.Write(buff)
}

func makeAGuess(w http.ResponseWriter, r *http.Request) {
	var guess userGuess

	params := mux.Vars(r)
	game, err := hangman.RetrieveGame(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Ready request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &guess)
	if err != nil {
		panic(err)
	}

	game = hangman.MakeAGuess(game, guess.Guess)
	hangman.UpdateGame(game)

	game, err = hangman.RetrieveGame(game.ID)
	responseJSON := gameInfoJSON{
		ID:             game.ID,
		TurnsLeft:      game.TurnsLeft,
		Used:           game.Used,
		AvailableHints: game.AvailableHints,
	}
	buff, error := json.MarshalIndent(responseJSON, "", "\t")
	if error != nil {
		log.Fatal("Could not serialize game")
	}

	w.Write(buff)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/games", newGame).Methods("GET")
	router.HandleFunc("/games/{id}", retrieveGameInfo).Methods("GET")
	router.HandleFunc("/games/{id}/guesses", makeAGuess).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8000", router))
}
