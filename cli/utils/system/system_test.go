package system

import (
	"testing"
	"time"
)

func TestFormatBytes(t *testing.T) {
	cases := []struct {
		input uint64
		want  string
	}{
		{0, "0 B"},
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1024 * 1024 * 5, "5.0 MB"},
	}

	for _, c := range cases {
		got := FormatBytes(c.input)
		if got != c.want {
			t.Errorf("FormatBytes(%d) = %q; want %q", c.input, got, c.want)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	d := 3661 * time.Second
	formatted := FormatDuration(d)
	if formatted != "1h 1m 1s" {
		t.Errorf("FormatDuration(%v) = %q; want '1h 1m 1s'", d, formatted)
	}
}

func TestMenuRuntime(t *testing.T) {
	res := MenuRuntime(3661)
	if res != "1 hr, 1 min, 1 sec" {
		t.Errorf("MenuRuntime(3661) = %q; want '1 hr, 1 min, 1 sec'", res)
	}
}
