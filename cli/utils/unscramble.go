package cliutils

import (
	"whatsrook/cli/utils/games"
)

// UnscrambleState defines the phase of an Unscramble match.
type UnscrambleState = games.UnscrambleState

const (
	UnscrambleStateLobby      = games.UnscrambleStateLobby
	UnscrambleStateInProgress = games.UnscrambleStateInProgress
)

// UnscramblePlayer represents a player in an Unscramble game.
type UnscramblePlayer = games.UnscramblePlayer

// UnscrambleGame represents a match session of the Unscramble game.
type UnscrambleGame = games.UnscrambleGame

var (
	// IsUnscrambleGameActive returns true if there is an active Unscramble game in the chat.
	IsUnscrambleGameActive = games.IsUnscrambleGameActive

	// GetUnscrambleGame returns the active game for a chat key, or nil.
	GetUnscrambleGame = games.GetUnscrambleGame

	// CreateUnscrambleGame creates a new lobby for a chat.
	CreateUnscrambleGame = games.CreateUnscrambleGame

	// DeleteUnscrambleGame removes a game from the active map.
	DeleteUnscrambleGame = games.DeleteUnscrambleGame

	// GetTurnTimeLimit calculates turn duration.
	GetTurnTimeLimit = games.GetTurnTimeLimit

	// GetRandomWord returns a random word of the given length and its scrambled version.
	GetRandomWord = games.GetRandomWord

	// GetCXPTitle maps cumulative XP (CXP) to player titles.
	GetCXPTitle = games.GetCXPTitle
)
