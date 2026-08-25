package utils

import (
	"whatsrook/utils/media"
)

// AudioPTTMeta contains converted Opus OGG data, duration in seconds, and 64-bin amplitude waveform bytes.
type AudioPTTMeta = media.AudioPTTMeta

var (
	// EnsureJPEG decodes input image bytes and re-encodes to valid JPEG bytes.
	EnsureJPEG = media.EnsureJPEG

	// EnsureOpusPTT converts audio bytes to WhatsApp Opus OGG format with waveform.
	EnsureOpusPTT = media.EnsureOpusPTT

	// ExtractWaveformForTest extracts duration and waveform for testing.
	ExtractWaveformForTest = media.ExtractWaveformForTest

	// TranscodeToMP3 converts input audio to MP3 via ffmpeg.
	TranscodeToMP3 = media.TranscodeToMP3

	// PrepareCallVideo converts video to audio (.mp3) and Annex-B H.264 stream (.h264).
	PrepareCallVideo = media.PrepareCallVideo

	// AudioDuration calculates MP3 duration in pure Go.
	AudioDuration = media.AudioDuration

	// SplitAnnexBAccessUnits splits raw H.264 stream into individual access units.
	SplitAnnexBAccessUnits = media.SplitAnnexBAccessUnits

	// AnnexBHasIDR checks if Annex-B video payload contains an IDR keyframe.
	AnnexBHasIDR = media.AnnexBHasIDR

	// ExtensionFor returns the file extension for a MIME type.
	ExtensionFor = media.ExtensionFor
)
