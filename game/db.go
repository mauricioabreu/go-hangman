package hangman

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	// Used to access pgsql driver
	_ "github.com/lib/pq"
)

// GameRow : game database type
type GameRow struct {
	ID             string
	TurnsLeft      int
	Used           string
	AvailableHints int
}

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

	log.Printf("ID of inserted game: %s", game.ID)
	return game.ID
}

// RetrieveGame : Retrieve a game from the database
func RetrieveGame(id string) (GameRow, error) {
	var game GameRow

	connStr := "user=postgres dbname=hangman password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	row := db.QueryRow("SELECT uuid, turns_left, used, available_hints FROM hangman.games WHERE uuid = $1", id)
	err = row.Scan(&game.ID, &game.TurnsLeft, &game.Used, &game.AvailableHints)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return GameRow{}, err
	case nil:
		fmt.Println(row)
		return game, nil
	default:
		panic(err)
	}
}

func toString(arr []string) string {
	return strings.Join(arr[:], "")
}
