package media

import (
	"testing"
	"time"

	"wa-core/proto/waE2E"
	"wa-core/types"
	"wa-core/types/events"
)

func TestRecordAndGetRecentMessage(t *testing.T) {
	msgID := "test-msg-12345"
	chatJID := types.NewJID("123456789", types.DefaultUserServer)
	senderJID := types.NewJID("123456789", types.DefaultUserServer)

	text := "Hello world"
	evt := &events.Message{
		Info: types.MessageInfo{
			ID:        msgID,
			Chat:      chatJID,
			Sender:    senderJID,
			PushName:  "Tester",
			Timestamp: time.Now(),
		},
		Message: &waE2E.Message{
			Conversation: &text,
		},
	}

	RecordRecentMessage(evt)

	entry, ok := GetRecentMessage(msgID)
	if !ok {
		t.Fatalf("expected to find recorded message %s", msgID)
	}

	if entry.ID != msgID {
		t.Errorf("expected ID %q, got %q", msgID, entry.ID)
	}
	if entry.Text != text {
		t.Errorf("expected Text %q, got %q", text, entry.Text)
	}
	if entry.PushName != "Tester" {
		t.Errorf("expected PushName 'Tester', got %q", entry.PushName)
	}
}
