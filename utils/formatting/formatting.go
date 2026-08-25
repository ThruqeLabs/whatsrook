package formatting

import (
	"fmt"
	"strings"
)

// Bold returns plain text without WhatsApp bold formatting symbols (*).
func Bold(text string) string {
	return text
}

// Boldf formats text according to format specifier without bold formatting symbols.
func Boldf(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

// Italic returns plain text without WhatsApp italic formatting symbols (_).
func Italic(text string) string {
	return text
}

// Italicf formats text according to format specifier without italic formatting symbols.
func Italicf(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

// Code returns plain text without WhatsApp inline code formatting symbols (`).
func Code(text string) string {
	return text
}

// Codef formats text according to format specifier without inline code formatting symbols.
func Codef(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

// CodeBlock returns plain text without WhatsApp code block formatting symbols (```).
func CodeBlock(code string, lang ...string) string {
	return code
}

// Strike returns plain text without WhatsApp strikethrough formatting symbols (~).
func Strike(text string) string {
	return text
}

// Strikef formats text according to format specifier without strikethrough formatting symbols.
func Strikef(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

// Quote returns plain text without WhatsApp quote formatting symbols (>).
func Quote(text string) string {
	return text
}

// Quotef formats text according to format specifier without quote formatting symbols.
func Quotef(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

// RemoveEmojis strips emoji runes from a string.
func RemoveEmojis(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if (r >= 0x1F000 && r <= 0x1F9FF) || (r >= 0x2600 && r <= 0x27BF) || (r >= 0x1FA00 && r <= 0x1FAFF) || (r >= 0x1F1E0 && r <= 0x1F1FF) {
			continue // skip emoji
		}
		sb.WriteRune(r)
	}
	return sb.String()
}

// SanitizeJID replaces special characters in a JID string for safe file/key representation.
func SanitizeJID(s string) string {
	return strings.NewReplacer("@", "_at_", ":", "_", ".", "_").Replace(s)
}

// IsKnownLanguageCode checks if a language code or name is valid for TTS.
func IsKnownLanguageCode(lang string) bool {
	clean := strings.ToLower(strings.TrimSpace(lang))
	if clean == "" {
		return false
	}
	// Common ISO language codes
	known := map[string]bool{
		"en": true, "es": true, "fr": true, "de": true, "it": true, "pt": true,
		"ru": true, "ja": true, "ko": true, "zh": true, "ar": true, "hi": true,
		"tr": true, "nl": true, "pl": true, "sv": true, "id": true, "th": true,
		"vi": true, "he": true, "uk": true, "cs": true, "el": true, "hu": true,
		"ro": true, "sk": true, "da": true, "fi": true, "no": true, "sw": true,
	}
	return known[clean] || len(clean) == 2 || len(clean) == 5
}
