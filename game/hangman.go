package hangman

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Game : gameplay state
type Game struct {
	ID             string   // Game identifier
	State          string   // Game state
	TurnsLeft      int      // Remaining attempts
	Letters        []string // Letters in the word
	Used           []string // Good guesses
	AvailableHints int      // Total of hints available
}

// PickWord : Randomly get a word from a set of words.
func PickWord(words []string) string {
	rand.Seed(time.Now().Unix())
	wordIndex := rand.Intn(len(words))
	return words[wordIndex]
}

func letterInWord(guess string, letters []string) bool {
	for _, letter := range letters {
		if guess == letter {
			return true
		}
	}
	return false
}

// RevealWord : reveal the word by checking if the guesses made
// are part of the choosen word.
func RevealWord(letters []string, used []string) string {
	revealedWord := ""

	for _, wordLetter := range letters {
		if letterInWord(wordLetter, used) {
			revealedWord += wordLetter
		} else {
			revealedWord += "_"
		}
	}

	return revealedWord
}

func hasWon(letters []string, used []string) bool {
	ocurrences := 0
	for _, letter := range letters {
		for _, goodGuess := range used {
			if letter == goodGuess {
				ocurrences++
			}
		}
	}
	return ocurrences >= len(letters)
}

// AskForHint : Allow player to ask for a hint
func AskForHint(game Game, letters []string, used []string) (Game, string) {
	var possibleHints []string

	// Check which letters can be given as a hint
	// that were not used yet. If no good guess was found,
	// indicate any letter of the word.
	if len(used) > 0 {
		for _, letter := range letters {
			for _, goodGuess := range used {
				if letter != goodGuess {
					possibleHints = append(possibleHints, letter)
				}
			}
		}
	} else {
		possibleHints = letters
	}
	rand.Seed(time.Now().Unix())
	hintIndex := rand.Intn(len(possibleHints))
	game.AvailableHints--
	return game, possibleHints[hintIndex]
}

// NewGame : Start a new game
func NewGame(turnsLeft int, word string) Game {
	letters := strings.Split(word, "")
	return Game{ID: uuid.New().String(), State: "initial", TurnsLeft: turnsLeft, Letters: letters, Used: []string{}, AvailableHints: 3}
}

// MakeAGuess : Process the player guess
func MakeAGuess(game Game, guess string) Game {
	if letterInWord(guess, game.Letters) {
		// If already guessed this letter...
		if letterInWord(guess, game.Used) == true {
			game.State = "alreadyGuessed"
		} else {
			game.Used = append(game.Used, guess)
			game.State = "goodGuess"
			if hasWon(game.Letters, game.Used) == true {
				game.State = "won"
			}
		}
	} else {
		game.TurnsLeft--
		game.State = "badGuess"
		if game.TurnsLeft == 0 {
			game.State = "lost"
		}
	}

	return game
}
