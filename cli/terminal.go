package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync/atomic"

	"github.com/mattn/go-isatty"
)

var isInteractiveSession atomic.Bool

// IsInteractiveSession returns whether the current process is running in an interactive wizard session.
func IsInteractiveSession() bool {
	return isInteractiveSession.Load()
}

// SetInteractiveSession marks the current session as interactive.
func SetInteractiveSession(v bool) {
	isInteractiveSession.Store(v)
}

// IsInteractiveTerminal checks whether standard input is attached to an interactive terminal/TTY.
func IsInteractiveTerminal() bool {
	fd := os.Stdin.Fd()
	return isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd)
}

// IsCIOrDaemon returns true if running in a continuous integration or non-interactive container environment.
func IsCIOrDaemon() bool {
	return os.Getenv("CI") != "" ||
		os.Getenv("CONTINUOUS_INTEGRATION") != "" ||
		os.Getenv("BUILD_NUMBER") != "" ||
		os.Getenv("DEBIAN_FRONTEND") == "noninteractive" ||
		os.Getenv("WHATSROOK_HEADLESS") == "true" ||
		os.Getenv("WHATSROOK_HEADLESS") == "1"
}

// SpawnTerminal attempts to spawn the current executable inside a newly allocated interactive terminal window.
// It returns true if a terminal was successfully spawned, in which case the calling process can exit cleanly.
func SpawnTerminal() (bool, error) {
	if os.Getenv("WHATSROOK_SPAWNED_TERMINAL") == "1" {
		return false, nil
	}
	if IsCIOrDaemon() {
		return false, nil
	}

	exePath, err := os.Executable()
	if err != nil {
		return false, err
	}

	switch runtime.GOOS {
	case "windows":
		// Launch a new Windows Command Prompt window executing whatsrook with --interactive
		cmd := exec.Command("cmd.exe", "/c", "start", "", exePath, "--interactive")
		cmd.Env = append(os.Environ(), "WHATSROOK_SPAWNED_TERMINAL=1")
		if err := cmd.Start(); err != nil {
			return false, err
		}
		return true, nil

	case "darwin":
		// macOS: Tell Terminal.app to open and run the executable interactively
		script := fmt.Sprintf(`tell application "Terminal" to do script %q`, exePath+" --interactive")
		cmd := exec.Command("osascript", "-e", script)
		cmd.Env = append(os.Environ(), "WHATSROOK_SPAWNED_TERMINAL=1")
		if err := cmd.Start(); err != nil {
			return false, err
		}
		return true, nil

	case "linux", "freebsd", "openbsd", "netbsd":
		// Only attempt GUI terminal spawning if an X11 / Wayland display is present
		if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
			return false, nil
		}

		terminals := []string{
			"x-terminal-emulator",
			"gnome-terminal",
			"konsole",
			"xfce4-terminal",
			"alacritty",
			"kitty",
			"wezterm",
			"foot",
			"xterm",
			"urxvt",
		}

		for _, t := range terminals {
			if path, err := exec.LookPath(t); err == nil {
				var cmd *exec.Cmd
				switch t {
				case "gnome-terminal", "xfce4-terminal":
					cmd = exec.Command(path, "--", exePath, "--interactive")
				default:
					cmd = exec.Command(path, "-e", exePath, "--interactive")
				}
				cmd.Env = append(os.Environ(), "WHATSROOK_SPAWNED_TERMINAL=1")
				if err := cmd.Start(); err == nil {
					return true, nil
				}
			}
		}
		return false, nil

	default:
		return false, nil
	}
}

// PromptExit keeps the terminal window open until the user presses Enter.
// This is essential on Windows when double-clicked, preventing the console window from closing instantly.
func PromptExit() {
	if !IsInteractiveSession() {
		return
	}
	if !IsInteractiveTerminal() {
		return
	}
	fmt.Print("\n[WhatsRook] Press Enter to exit...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}

// PromptExitWithReader is a testable variant of PromptExit.
func PromptExitWithReader(in io.Reader, out io.Writer) {
	_, _ = fmt.Fprint(out, "\n[WhatsRook] Press Enter to exit...")
	reader := bufio.NewReader(in)
	_, _ = reader.ReadString('\n')
}
