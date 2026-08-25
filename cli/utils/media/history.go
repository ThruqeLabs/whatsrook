package media

import (
	"sync"
	"time"

	"wa-core/types"
	"wa-core/types/events"

	"whatsrook/utils"
)

// RecentMessageEntry represents a tracked incoming or outgoing message for context recall.
type RecentMessageEntry struct {
	ID           string
	Chat         types.JID
	Sender       types.JID
	PushName     string
	Text         string
	Timestamp    time.Time
	HasQuoted    bool
	QuotedID     string
	QuotedSender types.JID
	QuotedName   string
	QuotedText   string
}

// MaxRecentMessages limits the number of tracked messages in RAM.
const MaxRecentMessages = 1000

var (
	// RecentMessagesMu protects the in-memory recent messages cache.
	RecentMessagesMu sync.RWMutex

	// RecentMessages maps message ID to message entry.
	RecentMessages = make(map[string]RecentMessageEntry)

	// RecentMsgOrder tracks insertion order for eviction.
	RecentMsgOrder []string
)

// RecordRecentMessage caches an incoming message in memory for context / quote lookups.
func RecordRecentMessage(evt *events.Message) {
	if evt == nil || evt.Message == nil || evt.Info.ID == "" {
		return
	}
	RecentMessagesMu.Lock()
	defer RecentMessagesMu.Unlock()

	sender := evt.Info.Sender
	pushName := evt.Info.PushName
	text := utils.ExtractMessageText(evt)

	entry := RecentMessageEntry{
		ID:        evt.Info.ID,
		Chat:      evt.Info.Chat,
		Sender:    sender,
		PushName:  pushName,
		Text:      text,
		Timestamp: evt.Info.Timestamp,
	}

	ci := utils.GetContextInfoFromProto(evt.Message)
	if ci != nil && ci.QuotedMessage != nil {
		entry.HasQuoted = true
		if ci.StanzaID != nil {
			entry.QuotedID = *ci.StanzaID
		}
		if ci.Participant != nil && *ci.Participant != "" {
			if pj, err := types.ParseJID(*ci.Participant); err == nil {
				entry.QuotedSender = pj
			}
		} else if ci.RemoteJID != nil && *ci.RemoteJID != "" {
			if pj, err := types.ParseJID(*ci.RemoteJID); err == nil {
				entry.QuotedSender = pj
			}
		}
		entry.QuotedText = utils.ExtractTextFromProto(ci.QuotedMessage)
	}

	if _, exists := RecentMessages[evt.Info.ID]; !exists {
		if len(RecentMsgOrder) >= MaxRecentMessages {
			oldest := RecentMsgOrder[0]
			RecentMsgOrder = RecentMsgOrder[1:]
			delete(RecentMessages, oldest)
		}
		RecentMsgOrder = append(RecentMsgOrder, evt.Info.ID)
	}
	RecentMessages[evt.Info.ID] = entry
}

// GetRecentMessage retrieves a cached recent message by its ID.
func GetRecentMessage(id string) (RecentMessageEntry, bool) {
	RecentMessagesMu.RLock()
	defer RecentMessagesMu.RUnlock()
	entry, ok := RecentMessages[id]
	return entry, ok
}
