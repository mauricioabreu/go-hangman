package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"go-hangman/db"
	hangman "go-hangman/game"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	// Used to access pgsql driver
	_ "github.com/lib/pq"
)

type gameInfoJSON struct {
	ID             string          `json:"id"`
	TurnsLeft      int             `json:"turns_left"`
	Used           map[string]bool `json:"used"`
	AvailableHints int             `json:"available_hints"`
}

type userGuess struct {
	Guess string
}

func newGame(w http.ResponseWriter, r *http.Request) {
	words := hangman.ReadWordsFromFile("words/words.txt")
	choosenWord := hangman.PickWord(words)
	game := hangman.NewGame(3, choosenWord)
	database.DbStore.CreateGame(game)
	w.Header().Set("Location", strings.Join([]string{r.Host, "games", game.ID}, "/"))
}

func retrieveGameInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	game, err := database.DbStore.RetrieveGame(params["id"])

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
	buff, error := json.MarshalIndent(responseJSON, "", "    ")
	if error != nil {
		log.Fatal("Could not serialize game")
	}

	w.Write(buff)
}

func makeAGuess(w http.ResponseWriter, r *http.Request) {
	var guess userGuess

	params := mux.Vars(r)
	game, err := database.DbStore.RetrieveGame(params["id"])

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
	database.DbStore.UpdateGame(game)

	game, err = database.DbStore.RetrieveGame(game.ID)
	responseJSON := gameInfoJSON{
		ID:             game.ID,
		TurnsLeft:      game.TurnsLeft,
		Used:           game.Used,
		AvailableHints: game.AvailableHints,
	}
	buff, error := json.MarshalIndent(responseJSON, "", "    ")
	if error != nil {
		log.Fatal("Could not serialize game")
	}

	w.Write(buff)
}

func main() {
	var dbUser string
	var dbName string
	var dbPassword string

	flag.StringVar(&dbUser, "db_user", "postgres", "Database user")
	flag.StringVar(&dbName, "db_name", "hangman", "Database name")
	flag.StringVar(&dbPassword, "db_password", "postgres", "Database password")

	flag.Parse()

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", dbUser, dbName, dbPassword)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	database.InitStore(&database.DB{DB: db})

	router := mux.NewRouter()
	// Register HTTP endpoints
	router.HandleFunc("/games", newGame).Methods("POST")
	router.HandleFunc("/games/{id}", retrieveGameInfo).Methods("GET")
	router.HandleFunc("/games/{id}/guesses", makeAGuess).Methods("PUT")
	// Set logger handler for the server
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Fatal(http.ListenAndServe(":8000", loggedRouter))
}
