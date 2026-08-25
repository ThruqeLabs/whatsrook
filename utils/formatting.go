package utils

import (
	"whatsrook/utils/formatting"
)

var (
	// Bold returns plain text without WhatsApp bold formatting symbols (*).
	Bold = formatting.Bold

	// Boldf formats text according to format specifier without bold formatting symbols.
	Boldf = formatting.Boldf

	// Italic returns plain text without WhatsApp italic formatting symbols (_).
	Italic = formatting.Italic

	// Italicf formats text according to format specifier without italic formatting symbols.
	Italicf = formatting.Italicf

	// Code returns plain text without WhatsApp inline code formatting symbols (`).
	Code = formatting.Code

	// Codef formats text according to format specifier without inline code formatting symbols.
	Codef = formatting.Codef

	// CodeBlock returns plain text without WhatsApp code block formatting symbols (```).
	CodeBlock = formatting.CodeBlock

	// Strike returns plain text without WhatsApp strikethrough formatting symbols (~).
	Strike = formatting.Strike

	// Strikef formats text according to format specifier without strikethrough formatting symbols.
	Strikef = formatting.Strikef

	// Quote returns plain text without WhatsApp quote formatting symbols (>).
	Quote = formatting.Quote

	// Quotef formats text according to format specifier without quote formatting symbols.
	Quotef = formatting.Quotef
)
