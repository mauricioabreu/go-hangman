package database

import (
	"database/sql"
	"log"
	"strings"

	"github.com/mauricioabreu/go-hangman/game"
)

// Store : Handle how games are changed and retrieve from the database
type Store interface {
	CreateGame(game hangman.Game) error
	UpdateGame(game hangman.Game) error
	RetrieveGame(id string) (hangman.Game, error)
	DeleteGame(id string) (bool, error)
}

// DB : Implement the Store interface
type DB struct {
	DB *sql.DB
}

// CreateGame : Insert a new game into the database
func (store *DB) CreateGame(game hangman.Game) error {
	_, err := store.DB.Exec("INSERT INTO hangman.games (uuid, turns_left, word, used, available_hints) VALUES ($1, $2, $3, $4, $5)",
		game.ID, game.TurnsLeft, toString(game.Letters), mapToString(game.Used), game.AvailableHints)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Game ID %s inserted", game.ID)
	return nil
}

// UpdateGame : Update game state
func (store *DB) UpdateGame(game hangman.Game) error {
	_, err := store.DB.Exec("UPDATE hangman.games SET turns_left = $1, used = $2, available_hints = $3 WHERE uuid = $4",
		game.TurnsLeft, mapToString(game.Used), game.AvailableHints, game.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Game ID %s updated", game.ID)
	return nil
}

// RetrieveGame : Retrieve a game from the database
func (store *DB) RetrieveGame(id string) (hangman.Game, error) {
	var (
		uuid           string
		turnsLeft      int
		word           string
		used           string
		availableHints int
	)

	row := store.DB.QueryRow("SELECT uuid, turns_left, word, used, available_hints FROM hangman.games WHERE uuid = $1", id)
	err := row.Scan(&uuid, &turnsLeft, &word, &used, &availableHints)

	switch err {
	case sql.ErrNoRows:
		log.Printf("No rows were returned for game ID: %s\n", id)
		return hangman.Game{}, err
	case nil:
		return hangman.Game{ID: uuid,
			TurnsLeft:      turnsLeft,
			Letters:        strings.Split(word, ""),
			Used:           stringToMap(used),
			AvailableHints: availableHints,
		}, nil
	default:
		panic(err)
	}
}

// DeleteGame : remove a game from the database
func (store *DB) DeleteGame(id string) (bool, error) {
	result, err := store.DB.Exec("DELETE FROM hangman.games WHERE uuid = $1", id)
	if err == nil {
		// Check if there was a game to delete
		count, err := result.RowsAffected()
		if err == nil {
			if count > 0 {
				return true, err
			}
			return false, err
		}
	}
	return false, err
}

// DbStore : database
var DbStore Store

// InitStore : initialize the storage backend
func InitStore(s Store) {
	DbStore = s
}

func toString(arr []string) string {
	return strings.Join(arr[:], "")
}

func mapToString(m map[string]bool) string {
	str := ""
	for key := range m {
		str += key
	}
	return str
}

func stringToMap(str string) map[string]bool {
	m := make(map[string]bool)
	for _, letter := range strings.Split(str, "") {
		m[letter] = true
	}
	return m
}
