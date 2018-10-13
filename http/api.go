package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	hangman "github.com/mauricioabreu/go-hangman/game"

	database "github.com/mauricioabreu/go-hangman/db"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	// Used to access pgsql driver
	_ "github.com/lib/pq"
)

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

type userGuess struct {
	Guess string
}

func customNewGame(wordsFile string) http.HandlerFunc {
	words := hangman.ReadWordsFromFile(wordsFile)

	return func(w http.ResponseWriter, r *http.Request) {
		choosenWord := hangman.PickWord(words)
		game := hangman.NewGame(3, choosenWord)
		database.DbStore.CreateGame(game)
		w.WriteHeader(http.StatusNoContent)
		w.Header().Set("Location", strings.Join([]string{r.Host, "games", game.ID}, "/"))
	}
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
	result, err := database.DbStore.DeleteGame(params["id"])
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
		dbUser     string
		dbName     string
		dbPassword string
		wordsFile  string
	)

	flag.StringVar(&dbUser, "db_user", "postgres", "Database user")
	flag.StringVar(&dbName, "db_name", "hangman", "Database name")
	flag.StringVar(&dbPassword, "db_password", "postgres", "Database password")
	flag.StringVar(&wordsFile, "words_file", "words/words.txt", "Words file")

	flag.Parse()

	// check if the words file is accessible
	if _, err := os.Stat(wordsFile); err != nil {
		log.Fatalf("Could not open the words file: %s\n", err)
	}

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
	router.Use(commonMiddleware)
	// Register HTTP endpoints
	router.HandleFunc("/games", customNewGame(wordsFile)).Methods("POST")
	router.HandleFunc("/games/{id}", retrieveGameInfo).Methods("GET")
	router.HandleFunc("/games/{id}/guesses", makeAGuess).Methods("PUT")
	router.HandleFunc("/games/{id}", deleteGame).Methods("DELETE")
	// Set logger handler for the server
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	log.Println("Starting HTTP server...")
	log.Fatal(http.ListenAndServe(":8000", loggedRouter))
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
