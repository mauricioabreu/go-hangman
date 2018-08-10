package cli

import (
	"bufio"
	"fmt"
	"go-hangman/game"
	"log"
	"os"
	"strings"
)

type game struct {
	turnsLeft int      // Remaining attempts
	letters   []string // letters in the word
	used      []string // Good guesses
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
	return game{turnsLeft: turnsLeft, letters: letters, used: []string{}}
}

// Play : play the game
func Play() {
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
		fmt.Println("Guess a letter for the word:")
		guess, error := reader.ReadString('\n')
		guess = strings.TrimSpace(guess)

		if error != nil {
			log.Fatal("Could not read from terminal")
			os.Exit(1)
		}

		fmt.Printf("Your guess was '%s'\n", guess)

		if hangman.LetterInWord(guess, game.letters) {
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

		if hangman.HasWon(game.letters, game.used) == true {
			fmt.Println("YOU WON!!!")
			os.Exit(0)
		}

		fmt.Println(hangman.RevealWord(game.letters, game.used))
	}
}