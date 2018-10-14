package datastore

import (
	"errors"
	"fmt"
	"log"

	hangman "github.com/mauricioabreu/go-hangman/game"
)

// memoryStore : implementation for in-memory Store
type memoryStore struct {
	m map[string]hangman.Game
}

// NewMemoryStore : Initialize an in-memory store
func NewMemoryStore() (Store, error) {
	return &memoryStore{m: make(map[string]hangman.Game)}, nil
}

// CreateGame : Insert a new game into the in-memory store
func (memStore memoryStore) CreateGame(game hangman.Game) error {
	memStore.m[game.ID] = game
	log.Printf("Game ID %s inserted", game.ID)
	return nil
}

// UpdateGame : Update game state
func (memStore memoryStore) UpdateGame(game hangman.Game) error {
	memStore.m[game.ID] = game
	log.Printf("Game ID %s updated", game.ID)
	return nil
}

// RetrieveGame : Retrieve a game from the database
func (memStore memoryStore) RetrieveGame(id string) (hangman.Game, error) {
	game, ok := memStore.m[id]
	if !ok {
		errmsg := fmt.Sprintf("No game was found for ID: %s\n", id)
		log.Print(errmsg)
		return hangman.Game{}, errors.New(errmsg)
	}
	return game, nil
}

// DeleteGame : remove a game from the database
func (memStore memoryStore) DeleteGame(id string) (bool, error) {
	_, ok := memStore.m[id]
	if !ok {
		return false, nil
	}
	delete(memStore.m, id)
	return true, nil
}
