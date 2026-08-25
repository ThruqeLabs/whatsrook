package formatting

import (
	"testing"
)

func TestRemoveEmojis(t *testing.T) {
	input := "Hello 🌍 World! 🚀🔥"
	expected := "Hello  World! "
	got := RemoveEmojis(input)
	if got != expected {
		t.Errorf("RemoveEmojis(%q) = %q; want %q", input, got, expected)
	}
}

func TestSanitizeJID(t *testing.T) {
	input := "123456789@s.whatsapp.net"
	expected := "123456789_at_s_whatsapp_net"
	got := SanitizeJID(input)
	if got != expected {
		t.Errorf("SanitizeJID(%q) = %q; want %q", input, got, expected)
	}
}

func TestIsKnownLanguageCode(t *testing.T) {
	if !IsKnownLanguageCode("en") {
		t.Errorf("expected 'en' to be known")
	}
	if !IsKnownLanguageCode("ES") {
		t.Errorf("expected 'ES' to be known")
	}
	if IsKnownLanguageCode("") {
		t.Errorf("expected '' not to be known")
	}
}
