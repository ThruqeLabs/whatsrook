package cliutils

import (
	"whatsrook/cli/utils/meta"
)

var (
	// RunCmd runs an arbitrary shell command and returns its combined stdout+stderr output.
	RunCmd = meta.RunCmd

	// GetChatType returns the chat type based on the JID's server suffix.
	GetChatType = meta.GetChatType

	// GetTerminalWidth returns the width of the terminal.
	GetTerminalWidth = meta.GetTerminalWidth

	// WrapText wraps text into lines of at most maxWidth characters without splitting words.
	WrapText = meta.WrapText
)
