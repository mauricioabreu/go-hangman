package main

import (
	"go-hangman/game"
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

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/games", newGame).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
