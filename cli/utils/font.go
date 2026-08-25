package cliutils

import (
	"whatsrook/cli/utils/font"
)

// FontEntry describes a supported font style.
type FontEntry = font.FontEntry

var (
	// FormatTextResponseRaw formats a text response with monospace/active font format, removing asterisks and emojis.
	FormatTextResponseRaw = font.FormatTextResponseRaw

	// SetFontStyle sets the active font style for text conversion.
	SetFontStyle = font.SetFontStyle

	// GetFontStyle returns the currently active font style name.
	GetFontStyle = font.GetFontStyle

	// ConvertFontStyle transforms the input string to the currently active font style.
	ConvertFontStyle = font.ConvertFontStyle

	// NormalizeFancyText converts Unicode fancy font characters back into ASCII.
	NormalizeFancyText = font.NormalizeFancyText

	// IndexedFonts lists all available indexed font styles.
	IndexedFonts = font.IndexedFonts
)
