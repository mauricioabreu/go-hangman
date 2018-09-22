package cli

import (
	"bufio"
	"fmt"
	"go-hangman/game"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
)

type game struct {
	turnsLeft      int      // Remaining attempts
	letters        []string // letters in the word
	used           []string // Good guesses
	availableHints int      // Total of hints available
}

func welcomePlayer() {
	asciiArt := `
	 _   _    _    _   _  ____ __  __    _    _   _ 
	| | | |  / \  | \ | |/ ___|  \/  |  / \  | \ | |
	| |_| | / _ \ |  \| | |  _| |\/| | / _ \ |  \| |
	|  _  |/ ___ \| |\  | |_| | |  | |/ ___ \| |\  |
	|_| |_/_/   \_\_| \_|\____|_|  |_/_/   \_\_| \_|
													
	`
	fmt.Println(asciiArt)
}

func initializeGame(turnsLeft int, word string) game {
	letters := strings.Split(word, "")
	return game{turnsLeft: turnsLeft, letters: letters, used: []string{}, availableHints: 3}
}

// Play : play the game
func Play() {
	// Colored messages
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	boldGreen := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow)

	words := []string{
		"apple",
		"banana",
		"orange",
	}
	welcomePlayer()
	choosenWord := hangman.PickWord(words)
	fmt.Println("Your word has", len(choosenWord), "letters")
	game := initializeGame(2, choosenWord)
	reader := bufio.NewReader(os.Stdin)
	// Game loop!
	for {
		fmt.Println("Guess a letter for the word or use '.h' for a hint:")
		guess, error := reader.ReadString('\n')
		guess = strings.TrimSpace(guess)

		if error != nil {
			log.Fatal("Could not read from terminal")
			os.Exit(1)
		}

		if guess == ".h" {
			if game.availableHints == 0 {
				red.Println("No more hints to use...")
				continue
			}
			guess = hangman.AskForHint(game.letters, game.used)
			game.availableHints--
		} else {
			fmt.Printf("Your guess was '%s'\n", guess)
		}

		if hangman.LetterInWord(guess, game.letters) {
			if hangman.LetterInWord(guess, game.used) == false {
				green.Println("Good guess!")
				game.used = append(game.used, guess)
			} else {
				yellow.Printf("Letter '%s' was already used...\n", guess)
			}
		} else {
			yellow.Printf("Sorry, '%s' is not in the word...\n", guess)
			game.turnsLeft--
		}

		if game.turnsLeft == 0 {
			red.Printf("You lost! The word was: %s\n", choosenWord)
			os.Exit(0)
		}

		if hangman.HasWon(game.letters, game.used) == true {
			boldGreen.Println("YOU WON!!!")
			green.Printf("The word was: %s\n", choosenWord)
			os.Exit(0)
		}

		fmt.Println(hangman.RevealWord(game.letters, game.used))
	}
}
