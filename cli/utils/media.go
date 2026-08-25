package cliutils

import (
	"whatsrook/cli/utils/media"
)

// RecentMessageEntry represents a tracked message entry.
type RecentMessageEntry = media.RecentMessageEntry

// MaxRecentMessages limits the number of tracked messages in RAM.
const MaxRecentMessages = media.MaxRecentMessages

var (
	// RecentMessagesMu protects the in-memory recent messages cache.
	RecentMessagesMu = &media.RecentMessagesMu

	// RecentMessages maps message ID to message entry.
	RecentMessages = media.RecentMessages

	// RecentMsgOrder tracks insertion order for eviction.
	RecentMsgOrder = media.RecentMsgOrder

	// RecordRecentMessage caches an incoming message in memory for context / quote lookups.
	RecordRecentMessage = media.RecordRecentMessage

	// GetRecentMessage retrieves a cached recent message by its ID.
	GetRecentMessage = media.GetRecentMessage
)
