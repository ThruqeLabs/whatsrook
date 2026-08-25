package utils

import (
	"whatsrook/utils/system"
)

// SystemMetadata contains runtime system and environment details.
type SystemMetadata = system.SystemMetadata

var (
	// GetGitCommit returns the short commit hash if running inside a Git repository.
	GetGitCommit = system.GetGitCommit

	// GetSystemMetadata gathers system metadata for diagnostics and status reporting.
	GetSystemMetadata = system.GetSystemMetadata

	// FormatBytes converts a byte count into a human-readable string (KB, MB, GB).
	FormatBytes = system.FormatBytes

	// FormatUptime converts seconds into a human-readable duration string.
	FormatUptime = system.FormatUptime
)
