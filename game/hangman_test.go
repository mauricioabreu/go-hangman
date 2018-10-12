package hangman

import (
	"strings"
	"testing"
)

func TestLetterInWord(t *testing.T) {
	word := []string{"f", "o", "o"}

	guess := "f"
	hasLetter := letterInWord(guess, word)
	if hasLetter != true {
		t.Errorf("Word %s does not contain letter %s", word, guess)
	}

	guess = "c"
	hasLetter = letterInWord(guess, word)
	if hasLetter == true {
		t.Errorf("Word %s should not contain letter %s", word, guess)
	}
}

func TestRevealWord(t *testing.T) {
	letters := []string{"f", "o", "o"}
	goodGuesses := map[string]bool{"o": true}

	revealedWord := RevealWord(letters, goodGuesses)
	if revealedWord != "_oo" {
		t.Errorf("It looks like the revealed word is not correct. Got %s, should be '_oo'", revealedWord)
	}

	letters = []string{"b", "a", "r"}
	goodGuesses = map[string]bool{"o": true}

	revealedWord = RevealWord(letters, goodGuesses)
	if revealedWord != "___" {
		t.Errorf("It looks like the revealed word is not correct. Got %s, should be '___'", revealedWord)
	}
}

func TestRevealWordSpecial(t *testing.T) {
	letters := []string{"s", "u", "g", "a", "r", "-", "f", "r", "e", "e"}
	goodGuesses := map[string]bool{"e": true}

	revealedWord := RevealWord(letters, goodGuesses)
	if revealedWord != "_____-__ee" {
		t.Errorf("It looks like the revealed word is not correct. Got %s, should be '_____-__ee'", revealedWord)
	}

	letters = []string{"b", "o", "b", "'", "s", " ", "p", "e", "n"}
	goodGuesses = map[string]bool{"b": true}

	revealedWord = RevealWord(letters, goodGuesses)
	if revealedWord != "b_b'_ ___" {
		t.Errorf("It looks like the revealed word is not correct. Got %s, should be 'b_b'_ ___", revealedWord)
	}
}

func TestWinsWithoutSpecial(t *testing.T) {
	g := NewGame(3, "big-al")
	g = MakeAGuess(g, "b")
	g = MakeAGuess(g, "i")
	g = MakeAGuess(g, "g")
	g = MakeAGuess(g, "a")
	g = MakeAGuess(g, "l")
	if g.State != "won" {
		t.Errorf("It looks like the game state is wrong. Got %s, should be 'won'", g.State)
	}

}

func TestAskForHint(t *testing.T) {
	g := NewGame(3, "love")
	availableHints := g.AvailableHints
	g, hint := AskForHint(g, g.Letters, make(map[string]bool))
	if strings.Contains("love", hint) != true {
		t.Errorf("The hint should be one of the letters in 'love'. Got %s", hint)
	}
	if availableHints-1 != g.AvailableHints {
		t.Errorf("After asking for a hint, we must decrement. Available hints so far: %d", g.AvailableHints)
	}
}

func TestAskForHintWithSpecial(t *testing.T) {
	g := NewGame(3, "sugar-free")
	g.GetRandomInt = func(i int) int { return 5 }

	availableHints := g.AvailableHints
	g, hint := AskForHint(g, g.Letters, make(map[string]bool))
	if strings.Contains("sugarfree", hint) != true {
		t.Errorf("The hint should be one of the letters in 'sugarfree'. Got %s", hint)
	}
	if availableHints-1 != g.AvailableHints {
		t.Errorf("After asking for a hint, we must decrement. Available hints so far: %d", g.AvailableHints)
	}
}
