package hangman

import "testing"

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

func TestRevealWordSpecial(t *testing.T) {
	letters := []string{"s", "u", "g", "a", "r", "-", "f", "r", "e", "e"}
	goodGuesses := []string{"e"}

	revealedWord := RevealWord(letters, goodGuesses)
	if revealedWord != "_____-__ee" {
		t.Errorf("It looks like the revealed word is not correct. Got %s, should be '_____-__ee'", revealedWord)
	}

	letters = []string{"b", "o", "b", "'", "s", " ", "p", "e", "n"}
	goodGuesses = []string{"b"}

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
