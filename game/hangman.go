package hangman

import (
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GetSystemRandomInt(i int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(i)
}

// Game : gameplay state
type Game struct {
	ID             string          // Game identifier
	State          string          // Game state
	TurnsLeft      int             // Remaining attempts
	Letters        []string        // Letters in the word
	Used           map[string]bool // Good guesses
	AvailableHints int             // Total of hints available
	GetRandomInt   func(int) int   // Source of randomness
}

// PickWord : Randomly get a word from a set of words.
func PickWord(words []string) string {
	wordIndex := GetSystemRandomInt(len(words))
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
// are part of the choosen word. Hyphens, apostrophies, and spaces are free.
func RevealWord(letters []string, used map[string]bool) string {
	revealedWord := ""

	for _, wordLetter := range letters {
		if used[wordLetter] {
			revealedWord += wordLetter
		} else if isSpecial(wordLetter) {
			revealedWord += wordLetter

		} else {
			revealedWord += "_"
		}
	}

	return revealedWord
}

func isSpecial(wordLetter string) bool {
	return strings.ContainsAny("-' ", wordLetter)
}

func hasWon(letters []string, used map[string]bool) bool {
	occurrences := 0
	for _, letter := range letters {
		if used[letter] || isSpecial(letter) {
			occurrences++
		}
	}
	return occurrences >= len(letters)
}

// AskForHint : Allow player to ask for a hint
func AskForHint(game Game, letters []string, used map[string]bool) (Game, string) {
	var validLetters, possibleHints []string

	// Filter out non-alphabetic characters from pool of hint
	// characters
	for _, letter := range letters {
		if "a" <= letter && letter <= "z" {
			validLetters = append(validLetters, letter)
		}
	}

	// Check which letters can be given as a hint
	// that were not used yet. If no good guess was found,
	// indicate any letter of the word.
	if len(used) > 0 {
		for _, letter := range validLetters {
			if !used[letter] {
				possibleHints = append(possibleHints, letter)
			}
		}
	} else {
		possibleHints = validLetters
	}

	hint := possibleHints[game.GetRandomInt(len(possibleHints))]
	game.State = "gotHint"
	game.Used[hint] = true
	game.AvailableHints--
	return game, hint
}

// NewGame : Start a new game
func NewGame(turnsLeft int, word string) Game {
	letters := strings.Split(word, "")
	return Game{ID: uuid.New().String(),
		State:          "initial",
		TurnsLeft:      turnsLeft,
		Letters:        letters,
		Used:           make(map[string]bool),
		AvailableHints: 3,
		GetRandomInt:   GetSystemRandomInt,
	}
}

// MakeAGuess : Process the player guess
func MakeAGuess(game Game, guess string) Game {
	if letterInWord(guess, game.Letters) {
		// If already guessed this letter...
		if game.Used[guess] {
			game.State = "alreadyGuessed"
		} else {
			game.Used[guess] = true
			game.State = "goodGuess"
			if hasWon(game.Letters, game.Used) {
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

// ReadWordsFromFile : Retrieve words from a text file
func ReadWordsFromFile(filePath string) []string {
	b, err := ioutil.ReadFile(filePath) // read words from file
	if err != nil {
		log.Fatal(err)
	}

	str := string(b) // convert content to a 'string'
	words := strings.Split(str, "\n")
	return words
}
