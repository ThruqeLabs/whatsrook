package games

import (
	"testing"
)

func TestGetRandomWord(t *testing.T) {
	orig, scram := GetRandomWord(5)
	if len(orig) != 5 {
		t.Errorf("expected 5-letter word, got %q (len %d)", orig, len(orig))
	}
	if len(scram) != 5 {
		t.Errorf("expected 5-letter scrambled word, got %q (len %d)", scram, len(scram))
	}
}

func TestGetTurnTimeLimit(t *testing.T) {
	if limit := GetTurnTimeLimit(3); limit != 30 {
		t.Errorf("expected 30s limit for level 3, got %d", limit)
	}
	if limit := GetTurnTimeLimit(16); limit != 6 {
		t.Errorf("expected 6s limit for level 16, got %d", limit)
	}
}

func TestGetRandomStartingLetter(t *testing.T) {
	r := GetRandomStartingLetter()
	if r < 'A' || r > 'Z' {
		t.Errorf("expected uppercase A-Z rune, got %c", r)
	}
}
