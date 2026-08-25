package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// TUI Styles and Color Palette
var (
	// Colors
	colorCyan      = lipgloss.Color("#00D9FF")
	colorGreen     = lipgloss.Color("#00FF88")
	colorPurple    = lipgloss.Color("#7952FF")
	colorYellow    = lipgloss.Color("#FFD000")
	colorRed       = lipgloss.Color("#FF4466")
	colorMuted     = lipgloss.Color("#6272A4")
	colorDarkBg    = lipgloss.Color("#13151F")
	colorSurfaceBg = lipgloss.Color("#1E2233")
	colorWhite     = lipgloss.Color("#F8F8F2")

	// Header Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#000000")).
			Background(colorCyan).
			Padding(0, 1)

	headerBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPurple).
			Padding(0, 1).
			Background(colorDarkBg)

	// Status Badges
	badgeConnected = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true)

	badgeConnecting = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true)

	badgeDisconnected = lipgloss.NewStyle().
				Foreground(colorRed).
				Bold(true)

	pillStyle = lipgloss.NewStyle().
			Foreground(colorWhite).
			Background(colorSurfaceBg).
			Padding(0, 1)

	// Box Containers
	logBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorCyan).
			Background(colorDarkBg)

	infoBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorMuted).
			Padding(0, 1).
			Background(colorDarkBg)

	modalBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(colorCyan).
			Padding(1, 2).
			Background(colorSurfaceBg)

	// Interactive Wizard Elements
	wizardTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorCyan).
			MarginBottom(1)

	stepNumber = lipgloss.NewStyle().
			Foreground(colorPurple).
			Bold(true)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(colorGreen).
				Bold(true)

	unselectedItemStyle = lipgloss.NewStyle().
				Foreground(colorWhite)

	cursorStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Bold(true)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Bold(true)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	promptStyle = lipgloss.NewStyle().
			Foreground(colorPurple).
			Bold(true)

	pairCodeBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(colorGreen).
				Padding(1, 3).
				Background(colorSurfaceBg).
				Bold(true)

	pairCodeNumberStyle = lipgloss.NewStyle().
				Foreground(colorGreen).
				Bold(true).
				Background(colorDarkBg).
				Padding(0, 2)

	notifModalBox = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(colorYellow).
			Padding(1, 3).
			Background(colorSurfaceBg).
			Align(lipgloss.Center)

	notifTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorYellow).
			MarginBottom(1)

	notifCountStyle = lipgloss.NewStyle().
			Foreground(colorCyan).
			Bold(true)

	btnActiveConfirm = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorGreen).
				Foreground(colorGreen).
				Background(colorDarkBg).
				Bold(true).
				Padding(0, 2)

	btnInactiveConfirm = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorMuted).
				Foreground(colorWhite).
				Padding(0, 2)

	btnActiveCancel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorRed).
			Foreground(colorRed).
			Background(colorDarkBg).
			Bold(true).
			Padding(0, 2)

	btnInactiveCancel = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorMuted).
				Foreground(colorMuted).
				Padding(0, 2)
)
