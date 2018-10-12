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

func initializeGame(turnsLeft int, word string) hangman.Game {
	return hangman.NewGame(turnsLeft, word)
}

// Play : play the game
func Play() {
	// Colored messages
	red := color.New(color.FgRed)
	green := color.New(color.FgGreen)
	boldGreen := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow)

	words := hangman.ReadWordsFromFile("words/words.txt")
	welcomePlayer()
	choosenWord := hangman.PickWord(words)
	fmt.Println("Your word has", len(choosenWord), "letters")
	game := initializeGame(3, choosenWord)
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
			if game.AvailableHints == 0 {
				red.Println("No more hints to use...")
				continue
			}
			game, guess = hangman.AskForHint(game, game.Letters, game.Used)
		} else {
			fmt.Printf("Your guess was '%s'\n", guess)
			game = hangman.MakeAGuess(game, guess)
		}

		if game.State == "goodGuess" {
			green.Println("Good guess!")
		}

		if game.State == "alreadyGuessed" {
			yellow.Printf("Letter '%s' was already used...\n", guess)
		}

		if game.State == "badGuess" {
			yellow.Printf("Sorry, '%s' is not in the word...\n", guess)
		}

		if game.State == "gotHint" {
			yellow.Printf("You have %d hints left.\n", game.AvailableHints)
		}

		if game.State == "lost" {
			red.Printf("You lost! The word was: %s\n", choosenWord)
			os.Exit(0)
		}

		if game.State == "won" {
			boldGreen.Println("YOU WON!!!")
			green.Printf("The word was: %s\n", choosenWord)
			os.Exit(0)
		}

		fmt.Println(hangman.RevealWord(game.Letters, game.Used))
	}
}
