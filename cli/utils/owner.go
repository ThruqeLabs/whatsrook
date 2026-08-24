package cliutils

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"

	"go.mau.fi/whatsmeow/types"
)

var (
	StatusBroadcastJID = types.JID{User: "status", Server: "broadcast"}

	ActiveShellSessions   = make(map[string]*ShellSession)
	ActiveShellSessionsMu sync.Mutex
	AnsiEscapeRegex       = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]|\x1b\([a-zA-Z]|\x1b\][0-9];[^\a\x1b]*(?:\a|\x1b\\)`)
)

type ShellSession struct {
	Chat           types.JID
	Sender         types.JID
	MsgID          types.MessageID
	Cmd            *exec.Cmd
	Stdin          io.WriteCloser
	Cancel         context.CancelFunc
	Buf            *bytes.Buffer
	Mu             sync.Mutex
	StartTime      time.Time
	CommandStr     string
	UpdateCh       chan struct{}
	Done           chan struct{}
	UserTerminated bool
}

// BuildShellCmd constructs an exec.Cmd appropriate for the host operating system.
func BuildShellCmd(ctx context.Context, commandStr string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		if _, err := exec.LookPath("pwsh.exe"); err == nil {
			return exec.CommandContext(ctx, "pwsh.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", commandStr)
		}
		if _, err := exec.LookPath("pwsh"); err == nil {
			return exec.CommandContext(ctx, "pwsh", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", commandStr)
		}
		if _, err := exec.LookPath("powershell.exe"); err == nil {
			return exec.CommandContext(ctx, "powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", commandStr)
		}
		if _, err := exec.LookPath("powershell"); err == nil {
			return exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", commandStr)
		}
		if _, err := exec.LookPath("cmd.exe"); err == nil {
			return exec.CommandContext(ctx, "cmd.exe", "/c", commandStr)
		}
		if _, err := exec.LookPath("cmd"); err == nil {
			return exec.CommandContext(ctx, "cmd", "/c", commandStr)
		}
		if _, err := exec.LookPath("bash.exe"); err == nil {
			return exec.CommandContext(ctx, "bash.exe", "-c", commandStr)
		}
		if _, err := exec.LookPath("bash"); err == nil {
			return exec.CommandContext(ctx, "bash", "-c", commandStr)
		}
		return exec.CommandContext(ctx, "powershell.exe", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command", commandStr)
	}

	shell := "bash"
	if _, err := exec.LookPath("bash"); err != nil {
		shell = "sh"
	}

	execCmdStr := commandStr
	if _, err := exec.LookPath("stdbuf"); err == nil {
		execCmdStr = "stdbuf -oL -eL " + commandStr
	}

	return exec.CommandContext(ctx, shell, "-c", execCmdStr)
}

func CleanShellOutput(raw string) string {
	cleaned := AnsiEscapeRegex.ReplaceAllString(raw, "")
	cleaned = strings.ReplaceAll(cleaned, "\r\n", "\n")
	lines := strings.Split(cleaned, "\n")
	var resultLines []string
	for _, line := range lines {
		if strings.Contains(line, "\r") {
			parts := strings.Split(line, "\r")
			var last string
			for _, part := range slices.Backward(parts) {
				trimmed := strings.TrimRight(part, " \t")
				if trimmed != "" {
					last = trimmed
					break
				}
			}
			if last != "" {
				resultLines = append(resultLines, last)
			}
		} else {
			resultLines = append(resultLines, line)
		}
	}
	return strings.Join(resultLines, "\n")
}
