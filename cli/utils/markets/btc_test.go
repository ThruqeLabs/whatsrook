package markets

import (
	"encoding/json"
	"strings"
	"testing"
)

const sampleBTCJSON = `{
  "meta": {
    "is_complete": false,
    "average_sec_between_blocks": 586.2,
    "previous_halving_block_number": 840000,
    "minus_blocks": 182300
  },
  "current": {
    "block_number": 867700,
    "timestamp": 1729864200
  },
  "target": {
    "block_number": 1050000,
    "predicted_timestamp": 1836864200
  },
  "bitcoin_price": {
    "price_usd": 68450.75,
    "time": 1729864200
  }
}`

func TestDecodeAndFormatBTC(t *testing.T) {
	var resp BTCPredictionsResponse
	err := json.Unmarshal([]byte(sampleBTCJSON), &resp)
	if err != nil {
		t.Fatalf("failed to unmarshal mock JSON: %v", err)
	}

	msg := FormatBTCMessage(&resp)
	if !strings.Contains(msg, "BITCOIN NETWORK & PRICE") {
		t.Errorf("expected header in formatted msg, got: %s", msg)
	}
	if !strings.Contains(msg, "$68450.75 USD") {
		t.Errorf("expected price $68450.75 USD in formatted msg, got: %s", msg)
	}
	if !strings.Contains(msg, "Current Block: 867700") {
		t.Errorf("expected current block in formatted msg, got: %s", msg)
	}
	if !strings.Contains(msg, "Blocks Remaining: 182300") {
		t.Errorf("expected blocks remaining in formatted msg, got: %s", msg)
	}
}
