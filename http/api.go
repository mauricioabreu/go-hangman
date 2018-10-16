package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mauricioabreu/go-hangman/datastore"
	hangman "github.com/mauricioabreu/go-hangman/game"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var store datastore.Store

// Retrieve all keys from a map
func keysFromMap(used map[string]bool) []string {
	keys := make([]string, 0, len(used))
	for k := range used {
		keys = append(keys, k)
	}
	return keys
}

type gameInfoJSON struct {
	ID             string   `json:"id"`
	Word           string   `json:"word"`
	TurnsLeft      int      `json:"turns_left"`
	Used           []string `json:"used"`
	AvailableHints int      `json:"available_hints"`
}

func (a *gameInfoJSON) equals(b gameInfoJSON) bool {
	if a.ID != b.ID || a.TurnsLeft != b.TurnsLeft || a.Word != b.Word {
		return false
	}
	used := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(b.Used)), ","), "[]")
	for _, ele := range a.Used {
		if !strings.Contains(used, ele) {
			return false
		}
	}
	return true
}

type userGuess struct {
	Guess string
}

func customNewGame(wordsFile string) http.HandlerFunc {
	words := hangman.ReadWordsFromFile(wordsFile)

	return func(w http.ResponseWriter, r *http.Request) {
		choosenWord := hangman.PickWord(words)
		game := hangman.NewGame(3, choosenWord)
		store.CreateGame(game)
		w.Header().Set("Location", strings.Join([]string{r.Host, "games", game.ID}, "/"))
		w.WriteHeader(http.StatusNoContent)
	}
}

func retrieveGameInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	game, err := store.RetrieveGame(params["id"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	responseJSON := gameInfoJSON{
		ID:             game.ID,
		Word:           hangman.RevealWord(game.Letters, game.Used),
		TurnsLeft:      game.TurnsLeft,
		Used:           keysFromMap(game.Used),
		AvailableHints: game.AvailableHints,
	}
	buff, err := json.MarshalIndent(responseJSON, "", "    ")
	if err != nil {
		log.Fatalf("Could not serialize game. Error: %s", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buff)
}

func makeAGuess(w http.ResponseWriter, r *http.Request) {
	var guess userGuess

	params := mux.Vars(r)
	game, err := store.RetrieveGame(params["id"])

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
	store.UpdateGame(game)

	game, err = store.RetrieveGame(game.ID)
	responseJSON := gameInfoJSON{
		ID:             game.ID,
		Word:           hangman.RevealWord(game.Letters, game.Used),
		TurnsLeft:      game.TurnsLeft,
		Used:           keysFromMap(game.Used),
		AvailableHints: game.AvailableHints,
	}
	buff, err := json.MarshalIndent(responseJSON, "", "    ")
	if err != nil {
		log.Fatalf("Could not serialize game. Error: %s", err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(buff)
}

func deleteGame(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	result, err := store.DeleteGame(params["id"])
	if err != nil {
		log.Fatalf("Could not delete the game. Error: %s", err)
	}
	if result {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func main() {
	var (
		storeType  string
		dbUser     string
		dbName     string
		dbPassword string
		wordsFile  string
	)

	flag.StringVar(&storeType, "datastore", "pg", "Data store type")
	flag.StringVar(&dbUser, "db_user", "postgres", "Database user")
	flag.StringVar(&dbName, "db_name", "hangman", "Database name")
	flag.StringVar(&dbPassword, "db_password", "postgres", "Database password")
	flag.StringVar(&wordsFile, "words_file", "words/words.txt", "Words file")

	flag.Parse()

	configureStoreType(storeType, dbName, dbUser, dbPassword)

	// check if the words file is accessible
	if _, err := os.Stat(wordsFile); err != nil {
		log.Fatalf("Could not open the words file: %s\n", err)
	}

	router := Router(wordsFile)
	// Set logger handler for the server
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Println("Starting HTTP server...")
	log.Fatal(http.ListenAndServe(":8000", loggedRouter))
}

func configureStoreType(storeType, dbName, dbUser, dbPassword string) {
	var err error
	switch storeType {
	case "pg":
		store, err = datastore.NewPgStore(dbName, dbUser, dbPassword)
	case "memory":
		store, err = datastore.NewMemoryStore()
	}
	if err != nil {
		log.Fatal(err)
	}
}

func Router(wordsFile string) *mux.Router {
	router := mux.NewRouter()
	router.Use(commonMiddleware)
	// Register HTTP endpoints
	router.HandleFunc("/games", customNewGame(wordsFile)).Methods("POST")
	router.HandleFunc("/games/{id}", retrieveGameInfo).Methods("GET")
	router.HandleFunc("/games/{id}/guesses", makeAGuess).Methods("PUT")
	router.HandleFunc("/games/{id}", deleteGame).Methods("DELETE")
	return router
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
