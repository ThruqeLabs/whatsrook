package utils

import (
	"io"
	Logger "whatsrook/logger"

	waLog "go.mau.fi/whatsmeow/util/log"
)

// InitLogger initializes the global logger with stdout and per-level log files in sessionDir/logs.
func InitLogger(sessionDir string, verbose bool) error {
	return Logger.InitLogger(sessionDir, verbose)
}

// InitLoggerWithOutput initializes the global logger with custom console output and per-level log files in sessionDir/logs.
func InitLoggerWithOutput(sessionDir string, verbose bool, consoleOut io.Writer) error {
	return Logger.InitLoggerWithOutput(sessionDir, verbose, consoleOut)
}

// CloseLogger flushes and closes all open log files.
func CloseLogger() {
	Logger.Close()
}

// WhatsrookLog creates a fast waLog.Logger adapter with module prefix.
func WhatsrookLog(module string, minLevel string, color bool) waLog.Logger {
	return Logger.WhatsrookLog(module, minLevel, color)
}
