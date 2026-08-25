package cliutils

import (
	"whatsrook/cli/utils/fun"
)

var (
	// GetRandomFact queries online random fact APIs with local fallback.
	GetRandomFact = fun.GetRandomFact

	// GetRandomQuote queries quote APIs with local fallback.
	GetRandomQuote = fun.GetRandomQuote

	// GetRandomJoke queries joke APIs with local fallback.
	GetRandomJoke = fun.GetRandomJoke

	// GetRandomRizz returns a random pickup line.
	GetRandomRizz = fun.GetRandomRizz
)
