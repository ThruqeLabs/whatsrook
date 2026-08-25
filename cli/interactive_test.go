package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCleanAndValidatePhone(t *testing.T) {
	cases := []struct {
		input    string
		expected string
		valid    bool
	}{
		{"+1 (234) 567-8901", "+12345678901", true},
		{"2348060598068", "2348060598068", true},
		{"+44 7911 123456", "+447911123456", true},
		{"12345", "", false},             // too short
		{"12345678901234567", "", false}, // too long (> 15 digits)
		{"abc-def", "", false},           // no digits
		{"", "", false},                  // empty
	}

	for _, c := range cases {
		got, ok := cleanAndValidatePhone(c.input)
		if ok != c.valid {
			t.Errorf("cleanAndValidatePhone(%q): expected valid=%v, got %v", c.input, c.valid, ok)
		}
		if ok && got != c.expected {
			t.Errorf("cleanAndValidatePhone(%q): expected %q, got %q", c.input, c.expected, got)
		}
	}
}

func TestRunInteractiveWizard_FullFlow(t *testing.T) {
	// Inputs for each step:
	// 1. Phone number: +1234567890
	// 2. Auth method: 1 (Pair code)
	// 3. Client identity: 1 (Chrome)
	// 4. Database: Enter (default sqlite)
	// 5. Redis: Enter (in-memory)
	// 6. Port: 3000
	// 7a. Skip old: Y
	// 7b. Verbose: N
	// 8. Save .env: n
	input := strings.Join([]string{
		"+1234567890",
		"1",
		"1",
		"",
		"",
		"3000",
		"y",
		"n",
		"n",
	}, "\n") + "\n"

	in := strings.NewReader(input)
	var out bytes.Buffer

	initial := CLIArgs{Port: 3000}
	res, proceed := RunInteractiveWizard(initial, in, &out)

	if !proceed {
		t.Fatalf("expected wizard to proceed, got proceed=false")
	}

	if res.Session != "+1234567890" {
		t.Errorf("expected session +1234567890, got %q", res.Session)
	}
	if !res.Pair || res.QRCode {
		t.Errorf("expected Pair=true, QRCode=false, got Pair=%v, QRCode=%v", res.Pair, res.QRCode)
	}
	if res.Client != "chrome" {
		t.Errorf("expected client chrome, got %q", res.Client)
	}
	if res.Database != "sqlite" {
		t.Errorf("expected database sqlite, got %q", res.Database)
	}
	if res.RedisURL != "" {
		t.Errorf("expected empty RedisURL, got %q", res.RedisURL)
	}
	if res.Port != 3000 {
		t.Errorf("expected port 3000, got %d", res.Port)
	}
	if !res.SkipOldMessages {
		t.Errorf("expected SkipOldMessages=true, got %v", res.SkipOldMessages)
	}
	if res.Verbose {
		t.Errorf("expected Verbose=false, got %v", res.Verbose)
	}

	output := out.String()
	if !strings.Contains(output, "Configuration Summary") {
		t.Errorf("expected output to contain Configuration Summary")
	}
}

func TestRunInteractiveWizard_StandbyFlow(t *testing.T) {
	// 1. Phone number: Enter (empty)
	// 1b. Confirm Standby: y
	// 2. Standby Port: 4000
	input := strings.Join([]string{
		"",
		"y",
		"4000",
	}, "\n") + "\n"

	in := strings.NewReader(input)
	var out bytes.Buffer

	initial := CLIArgs{Port: 3000}
	res, proceed := RunInteractiveWizard(initial, in, &out)

	if !proceed {
		t.Fatalf("expected wizard to proceed with standby mode")
	}

	if res.Session != "" {
		t.Errorf("expected empty session for standby mode, got %q", res.Session)
	}
	if res.Port != 4000 {
		t.Errorf("expected port 4000, got %d", res.Port)
	}

	output := out.String()
	if !strings.Contains(output, "Starting WhatsRook in Standby Mode") {
		t.Errorf("expected output to mention Standby Mode")
	}
}

func TestRunInteractiveWizard_CustomOptions(t *testing.T) {
	// Inputs:
	// 1. Invalid phone first, then valid phone: 123 -> +447911123456
	// 2. Auth: 2 (QR code)
	// 3. Client: 2 (Android)
	// 4. Database: postgres://postgres:secret@localhost:5432/whatsrook
	// 5. Redis: redis://localhost:6379/1
	// 6. Port: 8080
	// 7a. Skip old: n
	// 7b. Verbose: y
	// 8. Save .env: n
	input := strings.Join([]string{
		"123",
		"+447911123456",
		"2",
		"2",
		"postgres://postgres:secret@localhost:5432/whatsrook",
		"redis://localhost:6379/1",
		"8080",
		"n",
		"y",
		"n",
	}, "\n") + "\n"

	in := strings.NewReader(input)
	var out bytes.Buffer

	initial := CLIArgs{}
	res, proceed := RunInteractiveWizard(initial, in, &out)

	if !proceed {
		t.Fatalf("expected wizard to proceed")
	}

	if res.Session != "+447911123456" {
		t.Errorf("expected session +447911123456, got %q", res.Session)
	}
	if res.Pair || !res.QRCode {
		t.Errorf("expected Pair=false, QRCode=true")
	}
	if res.Client != "android" {
		t.Errorf("expected Client=android, got %q", res.Client)
	}
	if res.Database != "postgres://postgres:secret@localhost:5432/whatsrook" {
		t.Errorf("expected postgres db url, got %q", res.Database)
	}
	if res.RedisURL != "redis://localhost:6379/1" {
		t.Errorf("expected RedisURL redis://localhost:6379/1, got %q", res.RedisURL)
	}
	if res.Port != 8080 {
		t.Errorf("expected Port=8080, got %d", res.Port)
	}
	if res.SkipOldMessages {
		t.Errorf("expected SkipOldMessages=false, got %v", res.SkipOldMessages)
	}
	if !res.Verbose {
		t.Errorf("expected Verbose=true, got %v", res.Verbose)
	}
}

func TestRunInteractiveWizard_Cancel(t *testing.T) {
	// Empty reader (EOF / closed input immediately)
	in := strings.NewReader("")
	var out bytes.Buffer

	initial := CLIArgs{}
	_, proceed := RunInteractiveWizard(initial, in, &out)

	if proceed {
		t.Errorf("expected wizard to abort on EOF, got proceed=true")
	}
}

func TestSaveEnvConfig(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	args := CLIArgs{
		Session:         "+1234567890",
		Pair:            true,
		Client:          "ios",
		Database:        "postgres://user:pass@localhost:5432/db",
		RedisURL:        "redis://127.0.0.1:6379/0",
		Port:            3500,
		SkipOldMessages: true,
		Verbose:         false,
	}

	if err := SaveEnvConfig(args, envPath); err != nil {
		t.Fatalf("failed to save .env config: %v", err)
	}

	data, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("failed to read written .env: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "SESSION=+1234567890") {
		t.Errorf("expected .env to contain SESSION=+1234567890, got:\n%s", content)
	}
	if !strings.Contains(content, "PAIR=true") {
		t.Errorf("expected .env to contain PAIR=true, got:\n%s", content)
	}
	if !strings.Contains(content, "CLIENT=ios") {
		t.Errorf("expected .env to contain CLIENT=ios, got:\n%s", content)
	}
	if !strings.Contains(content, "PORT=3500") {
		t.Errorf("expected .env to contain PORT=3500, got:\n%s", content)
	}
	if !strings.Contains(content, "DATABASE_URL=postgres://user:pass@localhost:5432/db") {
		t.Errorf("expected .env to contain DATABASE_URL, got:\n%s", content)
	}

	// Now update .env with new values
	args2 := CLIArgs{
		Session:         "+447911123456",
		QRCode:          true,
		Client:          "android",
		Database:        "sqlite",
		Port:            3000,
		SkipOldMessages: false,
		Verbose:         true,
	}

	if err := SaveEnvConfig(args2, envPath); err != nil {
		t.Fatalf("failed to update .env config: %v", err)
	}

	data2, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("failed to read updated .env: %v", err)
	}
	content2 := string(data2)

	if !strings.Contains(content2, "SESSION=+447911123456") {
		t.Errorf("expected updated SESSION, got:\n%s", content2)
	}
	if !strings.Contains(content2, "QRCODE=true") {
		t.Errorf("expected QRCODE=true, got:\n%s", content2)
	}
	if !strings.Contains(content2, "CLIENT=android") {
		t.Errorf("expected CLIENT=android, got:\n%s", content2)
	}
	if !strings.Contains(content2, "VERBOSE=true") {
		t.Errorf("expected VERBOSE=true, got:\n%s", content2)
	}
}
