/*
Hangman is a very famous game that a player must guess
the word by suggesting letters being limited by a low number of
guesses.

More information here: https://en.wikipedia.org/wiki/Hangman_(game)
*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type game struct {
	turnsLeft int      // Remaining attempts
	letters   []string // letters in the word
	used      []string // Good guesses
}

// Randomly get a word from a set of words.
func pickWord(words []string) string {
	rand.Seed(time.Now().Unix())
	wordIndex := rand.Intn(len(words))
	return words[wordIndex]
}

// Check if word contains the given letter.
func letterInWord(guess string, letters []string) bool {
	for _, letter := range letters {
		if guess == letter {
			return true
		}
	}
	return false
}

func initializeGame(turnsLeft int, word string) game {
	letters := strings.Split(word, "")
	return game{turnsLeft: turnsLeft, letters: letters, used: []string{}}
}

// Reveal the word by checking if the guesses made
// are part of the choosen word.
func revealWord(letters []string, used []string) string {
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

// Check if the player has won the game
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

func main() {
	words := []string{
		"apple",
		"banana",
		"orange",
	}
	fmt.Println("Welcome to Hangman game")
	choosenWord := pickWord(words)
	fmt.Println("Your word has", len(choosenWord), "letters")
	game := initializeGame(2, choosenWord)

	// Game loop!
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Guess a letter for the word:")
		guess, error := reader.ReadString('\n')
		guess = strings.TrimSpace(guess)

		if error != nil {
			log.Fatal("Could not read from terminal")
			os.Exit(1)
		}

		fmt.Printf("Your guess was '%s'\n", guess)

		if letterInWord(guess, game.letters) {
			fmt.Println("Good guess!")
			game.used = append(game.used, guess)
		} else {
			fmt.Printf("Sorry, '%s' is not in the word...\n", guess)
			game.turnsLeft--
		}

		if game.turnsLeft == 0 {
			fmt.Printf("You lost! The word was: %s\n", choosenWord)
			os.Exit(0)
		}

		if hasWon(game.letters, game.used) == true {
			fmt.Println("YOU WON!!!")
			os.Exit(0)
		}

		fmt.Println(revealWord(game.letters, game.used))
	}
}
