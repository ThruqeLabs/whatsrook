package cliutils

import (
	"whatsrook/cli/utils/markets"
)

// BTCPredictionsResponse matches the response of https://api.watcher.guru/bitcoinhalving/predictions
type BTCPredictionsResponse = markets.BTCPredictionsResponse

var (
	// FetchBTCPredictions queries the Watcher Guru Bitcoin Halving and Price Predictions API.
	FetchBTCPredictions = markets.FetchBTCPredictions

	// FormatBTCMessage formats a BTCPredictionsResponse into a clean, human-readable WhatsApp message.
	FormatBTCMessage = markets.FormatBTCMessage
)
