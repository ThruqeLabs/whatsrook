package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestInteractiveSessionState(t *testing.T) {
	SetInteractiveSession(true)
	if !IsInteractiveSession() {
		t.Errorf("expected IsInteractiveSession() to be true")
	}
	SetInteractiveSession(false)
	if IsInteractiveSession() {
		t.Errorf("expected IsInteractiveSession() to be false")
	}
}

func TestIsCIOrDaemon(t *testing.T) {
	t.Setenv("CI", "true")
	if !IsCIOrDaemon() {
		t.Errorf("expected IsCIOrDaemon() to be true when CI=true")
	}

	t.Setenv("CI", "")
	t.Setenv("WHATSROOK_HEADLESS", "1")
	if !IsCIOrDaemon() {
		t.Errorf("expected IsCIOrDaemon() to be true when WHATSROOK_HEADLESS=1")
	}

	t.Setenv("WHATSROOK_HEADLESS", "")
	t.Setenv("DEBIAN_FRONTEND", "noninteractive")
	if !IsCIOrDaemon() {
		t.Errorf("expected IsCIOrDaemon() to be true when DEBIAN_FRONTEND=noninteractive")
	}
}

func TestPromptExitWithReader(t *testing.T) {
	in := strings.NewReader("\n")
	var out bytes.Buffer

	PromptExitWithReader(in, &out)

	output := out.String()
	if !strings.Contains(output, "Press Enter to exit...") {
		t.Errorf("expected prompt output to contain 'Press Enter to exit...', got %q", output)
	}
}
