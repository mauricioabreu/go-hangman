package hangman

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	// Used to access pgsql driver
	_ "github.com/lib/pq"
)

// CreateGame : Insert a new game into the database
func CreateGame(game Game) string {
	connStr := "user=postgres dbname=hangman password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO hangman.games (uuid, turns_left, word, used, available_hints) VALUES ($1, $2, $3, $4, $5)",
		game.ID, game.TurnsLeft, toString(game.Letters), toString(game.Used), game.AvailableHints)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Game ID %s inserted", game.ID)
	return game.ID
}

// UpdateGame : Update game state
func UpdateGame(game Game) {
	connStr := "user=postgres dbname=hangman password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("UPDATE hangman.games SET turns_left = $1, used = $2, available_hints = $3 WHERE uuid = $4",
		game.TurnsLeft, toString(game.Used), game.AvailableHints, game.ID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Game ID %s updated", game.ID)
}

// RetrieveGame : Retrieve a game from the database
func RetrieveGame(id string) (Game, error) {
	var (
		uuid           string
		turnsLeft      int
		word           string
		used           string
		availableHints int
	)

	connStr := "user=postgres dbname=hangman password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	row := db.QueryRow("SELECT uuid, turns_left, word, used, available_hints FROM hangman.games WHERE uuid = $1", id)
	err = row.Scan(&uuid, &turnsLeft, &word, &used, &availableHints)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return Game{}, err
	case nil:
		return Game{ID: uuid, TurnsLeft: turnsLeft, Letters: strings.Split(word, ""), Used: strings.Split(used, ""), AvailableHints: availableHints}, nil
	default:
		panic(err)
	}
}

func toString(arr []string) string {
	return strings.Join(arr[:], "")
}
