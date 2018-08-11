package hangman

import "testing"

func TestLetterInWord(t *testing.T) {
	word := []string{"f", "o", "o"}

	guess := "f"
	hasLetter := LetterInWord(guess, word)
	if hasLetter != true {
		t.Errorf("Word %s does not contain letter %s", word, guess)
	}

	guess = "c"
	hasLetter = LetterInWord(guess, word)
	if hasLetter == true {
		t.Errorf("Word %s should not contain letter %s", word, guess)
	}
}

func TestRevealWord(t *testing.T) {
	letters := []string{"f", "o", "o"}
	goodGuesses := []string{"o"}

	revealedWord := RevealWord(letters, goodGuesses)
	if revealedWord != "_oo" {
		t.Errorf("It looks like the revealed word is not correct. Got %s, should be '_oo'", revealedWord)
	}

	letters = []string{"b", "a", "r"}
	goodGuesses = []string{"o"}

	revealedWord = RevealWord(letters, goodGuesses)
	if revealedWord != "___" {
		t.Errorf("It looks like the revealed word is not correct. Got %s, should be '___'", revealedWord)
	}
}
