package markets

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// BTCPredictionsResponse matches the response of https://api.watcher.guru/bitcoinhalving/predictions
type BTCPredictionsResponse struct {
	Meta struct {
		IsComplete                 bool    `json:"is_complete"`
		AverageSecBetweenBlocks    float64 `json:"average_sec_between_blocks"`
		PreviousHalvingBlockNumber int64   `json:"previous_halving_block_number"`
		MinusBlocks                int64   `json:"minus_blocks"`
	} `json:"meta"`
	Current struct {
		BlockNumber int64 `json:"block_number"`
		Timestamp   int64 `json:"timestamp"`
	} `json:"current"`
	Target struct {
		BlockNumber        int64 `json:"block_number"`
		PredictedTimestamp int64 `json:"predicted_timestamp"`
	} `json:"target"`
	BitcoinPrice struct {
		PriceUSD float64 `json:"price_usd"`
		Time     int64   `json:"time"`
	} `json:"bitcoin_price"`
}

// FetchBTCPredictions queries the Watcher Guru Bitcoin Halving and Price Predictions API.
func FetchBTCPredictions(ctx context.Context) (*BTCPredictionsResponse, error) {
	apiURL := "https://api.watcher.guru/bitcoinhalving/predictions"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create BTC API request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to query BTC API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("BTC API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read BTC API response: %w", err)
	}

	var data BTCPredictionsResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to decode BTC API response: %w", err)
	}

	return &data, nil
}

// FormatBTCMessage formats a BTCPredictionsResponse into a clean, human-readable WhatsApp message.
func FormatBTCMessage(data *BTCPredictionsResponse, extra ...string) string {
	if data == nil {
		return ""
	}

	var sb strings.Builder

	// Header
	sb.WriteString("🪙 BITCOIN NETWORK & PRICE\n\n")

	// Current Price
	if data.BitcoinPrice.PriceUSD > 0 {
		sb.WriteString(fmt.Sprintf("💰 Price: $%.2f USD\n", data.BitcoinPrice.PriceUSD))
	}

	// Current Block & Time
	if data.Current.BlockNumber > 0 {
		sb.WriteString(fmt.Sprintf("📦 Current Block: %d\n", data.Current.BlockNumber))
	}
	if data.Current.Timestamp > 0 {
		t := time.Unix(data.Current.Timestamp, 0).UTC()
		sb.WriteString(fmt.Sprintf("🕒 Block Time: %s UTC\n", t.Format("2006-01-02 15:04:05")))
	}

	// Halving Prediction
	if data.Target.BlockNumber > 0 {
		sb.WriteString(fmt.Sprintf("\n🎯 Next Halving Target: Block %d\n", data.Target.BlockNumber))
	}
	if data.Meta.MinusBlocks > 0 {
		sb.WriteString(fmt.Sprintf("⏳ Blocks Remaining: %d\n", data.Meta.MinusBlocks))
	}
	if data.Target.PredictedTimestamp > 0 {
		predictedTime := time.Unix(data.Target.PredictedTimestamp, 0).UTC()
		sb.WriteString(fmt.Sprintf("📅 Estimated Halving Date: %s UTC\n", predictedTime.Format("2006-01-02 15:04:05")))
	}
	if data.Meta.AverageSecBetweenBlocks > 0 {
		sb.WriteString(fmt.Sprintf("⚡ Avg Block Time: %.1fs\n", data.Meta.AverageSecBetweenBlocks))
	}

	if len(extra) > 0 && extra[0] != "" {
		sb.WriteString("\n" + extra[0] + "\n")
	}

	return strings.TrimSpace(sb.String())
}
