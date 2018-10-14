package datastore

import (
	"database/sql"

	hangman "github.com/mauricioabreu/go-hangman/game"
)

// Store : Handle how games are changed and retrieve from the database
type Store interface {
	CreateGame(game hangman.Game) error
	UpdateGame(game hangman.Game) error
	RetrieveGame(id string) (hangman.Game, error)
	DeleteGame(id string) (bool, error)
}

// dbStore : Common implementation for database stores
type dbStore struct {
	DB *sql.DB
}
