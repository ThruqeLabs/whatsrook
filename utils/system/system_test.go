package system

import (
	"strings"
	"testing"
)

func TestFormatBytes(t *testing.T) {
	cases := []struct {
		input    uint64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1024 * 1024, "1.0 MB"},
		{1024 * 1024 * 1024 * 2, "2.0 GB"},
	}

	for _, c := range cases {
		got := FormatBytes(c.input)
		if got != c.expected {
			t.Errorf("FormatBytes(%d) = %q; want %q", c.input, got, c.expected)
		}
	}
}

func TestFormatUptime(t *testing.T) {
	cases := []struct {
		seconds  float64
		expected string
	}{
		{0, "0s"},
		{45, "45s"},
		{125, "2m 5s"},
		{3665, "1h 1m 5s"},
		{90065, "1d 1h 1m 5s"},
	}

	for _, c := range cases {
		got := FormatUptime(c.seconds)
		if got != c.expected {
			t.Errorf("FormatUptime(%f) = %q; want %q", c.seconds, got, c.expected)
		}
	}
}

func TestGetSystemMetadata(t *testing.T) {
	meta := GetSystemMetadata("1.0.0")
	if meta.Version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", meta.Version)
	}
	if meta.OS == "" || meta.Arch == "" || meta.NumCPU <= 0 || meta.GoVersion == "" {
		t.Errorf("incomplete metadata: %+v", meta)
	}
	if meta.Commit == "" {
		t.Errorf("empty commit metadata")
	}
	commit := GetGitCommit()
	if strings.TrimSpace(commit) == "" {
		t.Errorf("empty git commit")
	}
}
