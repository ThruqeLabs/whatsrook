package media

import (
	"testing"
)

func TestExtractWaveform(t *testing.T) {
	pcm := make([]byte, 16000) // 1 second of 8kHz 16-bit audio
	secs, waveform := ExtractWaveformForTest(pcm, 8000)
	if secs != 1 {
		t.Fatalf("expected 1 sec, got %d", secs)
	}
	if len(waveform) != 64 {
		t.Fatalf("expected 64 waveform bins, got %d", len(waveform))
	}
}

func TestExtensionFor(t *testing.T) {
	cases := map[string]string{
		"image/jpeg":             ".jpg",
		"image/png":              ".png",
		"image/webp":             ".webp",
		"audio/ogg; codecs=opus": ".ogg",
		"video/mp4":              ".mp4",
		"application/pdf":        ".pdf",
		"unknown/type":           ".bin",
	}
	for mime, expected := range cases {
		got := ExtensionFor(mime)
		if got != expected {
			t.Errorf("ExtensionFor(%q) = %q; want %q", mime, got, expected)
		}
	}
}
