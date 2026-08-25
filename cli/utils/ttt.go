package cliutils

import (
	"whatsrook/cli/utils/games"
)

// TTTGame models a Tic-Tac-Toe match state.
type TTTGame = games.TTTGame

var (
	// BotJID is a dummy JID representing the bot player in AI matches.
	BotJID = games.BotJID

	// WcgRng is the RNG for word chain game.
	WcgRng = games.WcgRng

	// GameRng is the shared RNG for games.
	GameRng = games.GameRng

	// GameHTTPClient is the HTTP client with timeouts for game APIs.
	GameHTTPClient = games.GameHTTPClient

	// TTTMu guards the active Tic-Tac-Toe games map.
	TTTMu = &games.TTTMu

	// TTTGames stores active Tic-Tac-Toe games per chat.
	TTTGames = games.TTTGames

	// WCGDictionary contains fallback dictionaries by word length.
	WCGDictionary = games.WCGDictionary

	// IsTTTGameActive returns true if an active game exists in the chat.
	IsTTTGameActive = games.IsTTTGameActive
)
