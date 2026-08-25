// Shell command execution helpers used internally by Meta AI tool invocations.
package meta

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"go.mau.fi/whatsmeow/types"
	"golang.org/x/term"
)

// RunCmd runs an arbitrary shell command and returns its combined stdout+stderr output.
func RunCmd(input string) (string, error) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// GetChatType returns the chat type based on the JID's server suffix.
func GetChatType(chatJID types.JID) string {
	switch chatJID.Server {
	case types.GroupServer:
		return "group"
	case types.NewsletterServer:
		return "newsletter"
	case types.BroadcastServer:
		return "broadcast"
	default:
		return "private"
	}
}

// GetTerminalWidth returns the width of the terminal.
func GetTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		return 80
	}
	return width
}

// WrapText wraps text into lines of at most maxWidth characters without splitting words.
func WrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		maxWidth = 80
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	var lines []string
	var current strings.Builder

	for _, word := range words {
		if current.Len() == 0 {
			current.WriteString(word)
		} else if current.Len()+1+len(word) <= maxWidth {
			current.WriteString(" ")
			current.WriteString(word)
		} else {
			lines = append(lines, current.String())
			current.Reset()
			current.WriteString(word)
		}
	}

	if current.Len() > 0 {
		lines = append(lines, current.String())
	}

	return lines
}
