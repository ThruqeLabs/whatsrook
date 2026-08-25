package cliutils

import (
	"context"
	"runtime"
	"strings"
	"testing"
)

func TestBuildShellCmd(t *testing.T) {
	ctx := context.Background()
	cmd := BuildShellCmd(ctx, "echo hello")
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}

	if runtime.GOOS == "windows" {
		path := strings.ToLower(cmd.Path)
		if !strings.Contains(path, "pwsh") && !strings.Contains(path, "powershell") && !strings.Contains(path, "cmd") && !strings.Contains(path, "bash") {
			t.Errorf("expected windows shell in cmd.Path, got %q", cmd.Path)
		}
	} else {
		path := strings.ToLower(cmd.Path)
		if !strings.Contains(path, "bash") && !strings.Contains(path, "sh") {
			t.Errorf("expected unix shell in cmd.Path, got %q", cmd.Path)
		}
	}
}

func TestCleanShellOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain text with crlf",
			input:    "line 1\r\nline 2\r\nline 3",
			expected: "line 1\nline 2\nline 3",
		},
		{
			name:     "ansi escape sequences",
			input:    "\x1b[32mSUCCESS\x1b[0m: Done\r\n",
			expected: "SUCCESS: Done\n",
		},
		{
			name:     "carriage return spinner overwrite",
			input:    "Loading 10%\rLoading 50%\rLoading 100%\nDone",
			expected: "Loading 100%\nDone",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CleanShellOutput(tt.input)
			if got != tt.expected {
				t.Errorf("CleanShellOutput() = %q, want %q", got, tt.expected)
			}
		})
	}
}
