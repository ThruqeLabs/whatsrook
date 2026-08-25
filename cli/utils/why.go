package cliutils

import (
	"whatsrook/cli/utils/fun"
)

// WhyRequest represents the payload sent to why.com's ultimate-search API.
type WhyRequest = fun.WhyRequest

// WhyPull represents a suggested follow-up question or related exploration topic.
type WhyPull = fun.WhyPull

// WhyResponse represents the response returned by why.com's ultimate-search API.
type WhyResponse = fun.WhyResponse

var (
	// DefaultWhyAPIURL is the official why.com API endpoint.
	DefaultWhyAPIURL = fun.DefaultWhyAPIURL

	// WhyAPIURL allows custom URL override (e.g. for testing).
	WhyAPIURL = fun.WhyAPIURL

	// QueryWhy performs a POST request to https://why.com/api/ultimate-search.
	QueryWhy = fun.QueryWhy
)
