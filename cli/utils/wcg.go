package cliutils

import (
	"whatsrook/cli/utils/games"
)

// WCGState defines the state of a Word Chain Game session.
type WCGState = games.WCGState

const (
	WCGStateLobby      = games.WCGStateLobby
	WCGStateInProgress = games.WCGStateInProgress
)

// WCGPlayer represents a player in a WCG match.
type WCGPlayer = games.WCGPlayer

// WCGGame represents a Word Chain Game session.
type WCGGame = games.WCGGame

var (
	// GetRandomStartingLetter returns a random uppercase letter from A to Z.
	GetRandomStartingLetter = games.GetRandomStartingLetter

	// IsWCGGameActive returns true if there is an active WCG game in the chat.
	IsWCGGameActive = games.IsWCGGameActive

	// GetWCGGame returns the active game for a chat key, or nil.
	GetWCGGame = games.GetWCGGame

	// CreateWCGGame creates a new lobby for a chat.
	CreateWCGGame = games.CreateWCGGame

	// DeleteWCGGame removes a game from the active map.
	DeleteWCGGame = games.DeleteWCGGame
)
