package hangman

import (
	"database/sql"
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

	log.Printf("ID of inserted game: %s", game.ID)
	return game.ID
}

func toString(arr []string) string {
	return strings.Join(arr[:], "")
}
