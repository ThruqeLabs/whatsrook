package cliutils

import (
	"whatsrook/cli/utils/markets"
)

// FFInstrumentResponse wraps the instrument query response.
type FFInstrumentResponse = markets.FFInstrumentResponse

// FFInstrumentData contains instrument metadata, metrics, and latest quotes.
type FFInstrumentData = markets.FFInstrumentData

// FFBarsResponse wraps candlestick bar metrics.
type FFBarsResponse = markets.FFBarsResponse

// FFListItem represents a single market instrument from the instrument list.
type FFListItem = markets.FFListItem

// FFListResponse wraps the instrument list query.
type FFListResponse = markets.FFListResponse

// SelectListRow is a single row for WhatsApp interactive list UI.
type SelectListRow = markets.SelectListRow

// SelectListSection is a section group for WhatsApp interactive list UI.
type SelectListSection = markets.SelectListSection

// SelectListParams contains list parameters for market instrument selection.
type SelectListParams = markets.SelectListParams

var (
	// NormalizeMarketPair normalizes ticker symbols and shorthand names to standard instrument pairs.
	NormalizeMarketPair = markets.NormalizeMarketPair

	// FetchForexFactoryInstrumentList fetches available market pairs.
	FetchForexFactoryInstrumentList = markets.FetchForexFactoryInstrumentList

	// FetchMarketBars retrieves candlestick chart bars for an instrument.
	FetchMarketBars = markets.FetchMarketBars

	// FetchSingleMarket retrieves real-time quote metrics for an instrument.
	FetchSingleMarket = markets.FetchSingleMarket

	// FetchAllMarkets retrieves real-time quote metrics for multiple instruments in a single query.
	FetchAllMarkets = markets.FetchAllMarkets
)
