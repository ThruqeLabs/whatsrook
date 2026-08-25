package timezone

import (
	"testing"
)

func TestResolveTimezoneAlias(t *testing.T) {
	cases := []struct {
		input    string
		expected string
		found    bool
	}{
		{"Europe/London", "Europe/London", true},
		{"Etc/GMT", "Etc/GMT", true},
		{"Pacific Standard Time", "America/Los_Angeles", true},
		{"NonExistentZoneXYZ", "", false},
	}

	for _, c := range cases {
		got, ok := ResolveTimezoneAlias(c.input)
		if ok != c.found {
			t.Errorf("ResolveTimezoneAlias(%q) found = %v; want %v", c.input, ok, c.found)
		}
		if ok && got != c.expected {
			t.Errorf("ResolveTimezoneAlias(%q) = %q; want %q", c.input, got, c.expected)
		}
	}
}
