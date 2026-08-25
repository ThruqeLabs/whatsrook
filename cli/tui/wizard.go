package tui

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Config holds configuration parameters passed into the TUI wizard and dashboard.
type Config struct {
	Session         string
	Pair            bool
	QRCode          bool
	Client          string
	Database        string
	RedisURL        string
	Port            int
	SkipOldMessages bool
	Verbose         bool
	Interactive     bool
	Idle            bool
	Logout          bool
	NoTUI           bool
	TUI             bool
}

type wizardStep int

const (
	stepExistingEnvMenu wizardStep = iota
	stepPhone
	stepAuthMethod
	stepClient
	stepDatabase
	stepRedis
	stepPort
	stepSkipOld
	stepVerbose
	stepSaveEnv
	stepSummary
)

// WizardModel is the Bubbletea Model for the interactive setup wizard.
type WizardModel struct {
	step          wizardStep
	args          Config
	initialArgs   Config
	existingEnv   bool
	envSession    string
	menuIndex     int
	choiceIndex   int
	summaryChoice int
	textInput     textinput.Model
	errMessage    string
	width         int
	height        int
	quitting      bool
	proceed       bool
	saveEnv       bool
}

// NewWizardModel creates a new wizard Bubbletea model.
func NewWizardModel(initial Config) WizardModel {
	ti := textinput.New()
	ti.Focus()
	ti.Prompt = "> "
	ti.PromptStyle = cursorStyle
	ti.TextStyle = selectedItemStyle
	ti.PlaceholderStyle = helpDescStyle
	ti.Cursor.SetChar("|")
	ti.Cursor.Style = cursorStyle
	ti.CharLimit = 64
	ti.Width = 35

	hasExisting := false
	existingSession := ""

	// Check if .env or initial args contains an existing session
	if initial.Session != "" {
		hasExisting = true
		existingSession = initial.Session
	} else if data, err := os.ReadFile(".env"); err == nil {
		lines := strings.SplitSeq(string(data), "\n")
		for l := range lines {
			l = strings.TrimSpace(l)
			if after, ok := strings.CutPrefix(l, "SESSION="); ok {
				val := after
				val = strings.Trim(val, `"' `)
				if val != "" {
					hasExisting = true
					existingSession = val
					initial.Session = val
					break
				}
			}
		}
	}

	startStep := stepPhone
	if hasExisting && !initial.Interactive {
		startStep = stepExistingEnvMenu
	}

	m := WizardModel{
		step:          startStep,
		args:          initial,
		initialArgs:   initial,
		existingEnv:   hasExisting,
		envSession:    existingSession,
		summaryChoice: 1, // Default to Confirm & Start
		textInput:     ti,
		saveEnv:       true,
	}

	m.setupStepInput()
	return m
}

func (m *WizardModel) setupStepInput() {
	m.errMessage = ""
	m.choiceIndex = 0

	switch m.step {
	case stepExistingEnvMenu:
		m.menuIndex = 0

	case stepPhone:
		m.textInput.Placeholder = "+1234567890 (or empty for Standby)"
		m.textInput.SetValue(m.args.Session)
		m.textInput.CursorEnd()

	case stepAuthMethod:
		if m.args.QRCode {
			m.choiceIndex = 1
		} else {
			m.choiceIndex = 0
		}

	case stepClient:
		switch strings.ToLower(m.args.Client) {
		case "android":
			m.choiceIndex = 1
		case "ios":
			m.choiceIndex = 2
		default:
			m.choiceIndex = 0
		}

	case stepDatabase:
		m.textInput.Placeholder = "sqlite (or postgres://...)"
		if m.args.Database != "" && m.args.Database != "sqlite" {
			m.textInput.SetValue(m.args.Database)
		} else {
			m.textInput.SetValue("")
		}
		m.textInput.CursorEnd()

	case stepRedis:
		m.textInput.Placeholder = "in-memory (or redis://...)"
		m.textInput.SetValue(m.args.RedisURL)
		m.textInput.CursorEnd()

	case stepPort:
		m.textInput.Placeholder = "3000"
		if m.args.Port > 0 {
			m.textInput.SetValue(strconv.Itoa(m.args.Port))
		} else {
			m.textInput.SetValue("3000")
		}
		m.textInput.CursorEnd()

	case stepSkipOld:
		if m.args.SkipOldMessages {
			m.choiceIndex = 0
		} else {
			m.choiceIndex = 1
		}

	case stepVerbose:
		if m.args.Verbose {
			m.choiceIndex = 0
		} else {
			m.choiceIndex = 1
		}

	case stepSaveEnv:
		m.choiceIndex = 0

	case stepSummary:
		m.summaryChoice = 1
	}
}

func (m WizardModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m WizardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			m.proceed = false
			return m, tea.Quit

		case "enter":
			return m.handleEnter()

		case "up", "k":
			if m.step == stepExistingEnvMenu {
				if m.menuIndex > 0 {
					m.menuIndex--
				}
			} else if m.isChoiceStep() {
				if m.choiceIndex > 0 {
					m.choiceIndex--
				}
			} else if m.step == stepSummary {
				m.summaryChoice = 0
			}

		case "down", "j":
			if m.step == stepExistingEnvMenu {
				if m.menuIndex < 3 {
					m.menuIndex++
				}
			} else if m.isChoiceStep() {
				maxChoice := m.maxChoices()
				if m.choiceIndex < maxChoice {
					m.choiceIndex++
				}
			} else if m.step == stepSummary {
				m.summaryChoice = 1
			}

		case "left", "h":
			if m.step == stepSummary {
				m.summaryChoice = 0 // Cancel button on left
			}

		case "right", "l":
			if m.step == stepSummary {
				m.summaryChoice = 1 // Confirm button on right
			}

		case "tab":
			if m.isChoiceStep() {
				m.choiceIndex = (m.choiceIndex + 1) % (m.maxChoices() + 1)
			} else if m.step == stepSummary {
				m.summaryChoice = 1 - m.summaryChoice
			}

		case "shift+tab":
			if m.isChoiceStep() {
				if m.choiceIndex > 0 {
					m.choiceIndex--
				} else {
					m.choiceIndex = m.maxChoices()
				}
			} else if m.step == stepSummary {
				m.summaryChoice = 1 - m.summaryChoice
			}
		}
	}

	var cmd tea.Cmd
	if !m.isChoiceStep() && m.step != stepExistingEnvMenu && m.step != stepSummary {
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

func (m WizardModel) isChoiceStep() bool {
	return m.step == stepAuthMethod ||
		m.step == stepClient ||
		m.step == stepSkipOld ||
		m.step == stepVerbose ||
		m.step == stepSaveEnv
}

func (m WizardModel) maxChoices() int {
	switch m.step {
	case stepAuthMethod, stepSkipOld, stepVerbose, stepSaveEnv:
		return 1 // 2 choices: 0 or 1
	case stepClient:
		return 2 // 3 choices: 0, 1, 2
	default:
		return 0
	}
}

func (m WizardModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case stepExistingEnvMenu:
		switch m.menuIndex {
		case 0: // Quick Start
			m.proceed = true
			return m, tea.Quit
		case 1: // Modify
			m.step = stepPhone
			m.setupStepInput()
			return m, nil
		case 2: // New Session
			m.args.Session = ""
			m.step = stepPhone
			m.setupStepInput()
			return m, nil
		case 3: // Standby Mode
			m.args.Session = ""
			m.args.Idle = true
			m.proceed = true
			return m, tea.Quit
		}

	case stepPhone:
		val := strings.TrimSpace(m.textInput.Value())
		if val == "" {
			// Standby mode
			m.args.Session = ""
			m.step = stepPort
			m.setupStepInput()
			return m, nil
		}
		cleaned, ok := CleanAndValidatePhone(val)
		if !ok {
			m.errMessage = "Invalid phone number. Please enter 7-15 digits with country code."
			return m, nil
		}
		m.args.Session = cleaned
		m.step = stepAuthMethod
		m.setupStepInput()
		return m, nil

	case stepAuthMethod:
		if m.choiceIndex == 0 {
			m.args.Pair = true
			m.args.QRCode = false
		} else {
			m.args.Pair = false
			m.args.QRCode = true
		}
		m.step = stepClient
		m.setupStepInput()
		return m, nil

	case stepClient:
		switch m.choiceIndex {
		case 0:
			m.args.Client = "chrome"
		case 1:
			m.args.Client = "android"
		case 2:
			m.args.Client = "ios"
		}
		m.step = stepDatabase
		m.setupStepInput()
		return m, nil

	case stepDatabase:
		val := strings.TrimSpace(m.textInput.Value())
		if val == "" || strings.EqualFold(val, "sqlite") {
			m.args.Database = "sqlite"
		} else {
			m.args.Database = val
		}
		m.step = stepRedis
		m.setupStepInput()
		return m, nil

	case stepRedis:
		val := strings.TrimSpace(m.textInput.Value())
		if val == "" || strings.EqualFold(val, "skip") || strings.EqualFold(val, "none") || strings.EqualFold(val, "in-memory") {
			m.args.RedisURL = ""
		} else {
			m.args.RedisURL = val
		}
		m.step = stepPort
		m.setupStepInput()
		return m, nil

	case stepPort:
		val := strings.TrimSpace(m.textInput.Value())
		if val == "" {
			if m.args.Port <= 0 {
				m.args.Port = 3000
			}
		} else {
			p, err := strconv.Atoi(val)
			if err != nil || p <= 0 || p > 65535 {
				m.errMessage = "Please enter a valid port between 1 and 65535."
				return m, nil
			}
			m.args.Port = p
		}
		if m.args.Session == "" {
			// In standby mode, jump straight to summary
			m.step = stepSummary
			return m, nil
		}
		m.step = stepSkipOld
		m.setupStepInput()
		return m, nil

	case stepSkipOld:
		m.args.SkipOldMessages = (m.choiceIndex == 0)
		m.step = stepVerbose
		m.setupStepInput()
		return m, nil

	case stepVerbose:
		m.args.Verbose = (m.choiceIndex == 0)
		m.step = stepSaveEnv
		m.setupStepInput()
		return m, nil

	case stepSaveEnv:
		m.saveEnv = (m.choiceIndex == 0)
		if m.saveEnv {
			_ = SaveEnvConfig(m.args, ".env")
		}
		m.step = stepSummary
		return m, nil

	case stepSummary:
		if m.summaryChoice == 0 {
			m.quitting = true
			m.proceed = false
			return m, tea.Quit
		}
		m.proceed = true
		return m, tea.Quit
	}

	return m, nil
}

func (m WizardModel) View() string {
	if m.quitting {
		return "\n  Setup cancelled.\n\n"
	}

	var sb strings.Builder

	// Top Title
	sb.WriteString(titleStyle.Render(" WhatsRook Setup Wizard ") + "\n\n")

	switch m.step {
	case stepExistingEnvMenu:
		sb.WriteString(wizardTitle.Render("Existing Session Detected") + "\n")

		cardContent := fmt.Sprintf(
			"Phone / Session:   %s\nClient Platform:   %s\nDatabase Storage:  %s\nServer Port:       %d",
			badgeConnected.Render(m.args.Session),
			pillStyle.Render(m.args.Client),
			pillStyle.Render(m.args.Database),
			m.args.Port,
		)
		sb.WriteString(infoBoxStyle.Render(cardContent) + "\n\n")

		menuItems := []string{
			"Quick Start (Launch with current settings)",
			"Modify / Reconfigure settings",
			"Add / Switch to a different session phone number",
			"Run in Standby (Idle) Mode",
		}

		for i, item := range menuItems {
			cursor := "  "
			itemText := unselectedItemStyle.Render(item)
			if i == m.menuIndex {
				cursor = cursorStyle.Render("> ")
				itemText = selectedItemStyle.Render(item)
			}
			sb.WriteString(fmt.Sprintf("%s%s\n", cursor, itemText))
		}
		sb.WriteString("\n" + helpDescStyle.Render("Up/Down: navigate | Enter: select | Esc: quit"))

	case stepPhone:
		sb.WriteString(stepNumber.Render("[Step 1/8] ") + wizardTitle.Render("WhatsApp Phone Number") + "\n")
		sb.WriteString("Enter your full WhatsApp phone number with country code.\n")
		sb.WriteString(helpDescStyle.Render("Leave blank and press Enter to run in Standby Mode.\n\n"))
		sb.WriteString(promptStyle.Render("Phone: ") + m.textInput.View() + "\n")
		if m.errMessage != "" {
			sb.WriteString("\n" + badgeDisconnected.Render("! "+m.errMessage) + "\n")
		}
		sb.WriteString("\n" + helpDescStyle.Render("Enter: continue | Esc: cancel"))

	case stepAuthMethod:
		sb.WriteString(stepNumber.Render("[Step 2/8] ") + wizardTitle.Render("Linking & Authentication Method") + "\n")
		sb.WriteString("Choose how you want to link your WhatsApp account:\n\n")

		options := []string{
			"Pair Code (Recommended - 8-character code sent to your WhatsApp app)",
			"QR Code (Scan QR code via browser or terminal)",
		}
		for i, opt := range options {
			radio := "( ) "
			text := unselectedItemStyle.Render(opt)
			if i == m.choiceIndex {
				radio = cursorStyle.Render("(*) ")
				text = selectedItemStyle.Render(opt)
			}
			sb.WriteString(fmt.Sprintf("  %s%s\n", radio, text))
		}
		sb.WriteString("\n" + helpDescStyle.Render("Up/Down: navigate | Enter: select | Esc: cancel"))

	case stepClient:
		sb.WriteString(stepNumber.Render("[Step 3/8] ") + wizardTitle.Render("Client Device Identity") + "\n")
		sb.WriteString("Select the device platform that WhatsRook will emulate:\n\n")

		options := []string{
			"Chrome / Web (Recommended - default)",
			"Android Phone",
			"iOS Phone (iPhone)",
		}
		for i, opt := range options {
			radio := "( ) "
			text := unselectedItemStyle.Render(opt)
			if i == m.choiceIndex {
				radio = cursorStyle.Render("(*) ")
				text = selectedItemStyle.Render(opt)
			}
			sb.WriteString(fmt.Sprintf("  %s%s\n", radio, text))
		}
		sb.WriteString("\n" + helpDescStyle.Render("Up/Down: navigate | Enter: select | Esc: cancel"))

	case stepDatabase:
		sb.WriteString(stepNumber.Render("[Step 4/8] ") + wizardTitle.Render("Database Storage") + "\n")
		sb.WriteString("WhatsRook uses SQLite by default (local whatsrook.db - zero configuration).\n")
		sb.WriteString(helpDescStyle.Render("Press Enter to use SQLite, or enter a PostgreSQL connection URL:\n\n"))
		sb.WriteString(promptStyle.Render("Database: ") + m.textInput.View() + "\n")
		sb.WriteString("\n" + helpDescStyle.Render("Enter: continue | Esc: cancel"))

	case stepRedis:
		sb.WriteString(stepNumber.Render("[Step 5/8] ") + wizardTitle.Render("Cache Backend (Optional)") + "\n")
		sb.WriteString("WhatsRook uses fast in-memory caching by default.\n")
		sb.WriteString(helpDescStyle.Render("Press Enter to use in-memory cache, or enter a Redis URL:\n\n"))
		sb.WriteString(promptStyle.Render("Redis URL: ") + m.textInput.View() + "\n")
		sb.WriteString("\n" + helpDescStyle.Render("Enter: continue | Esc: cancel"))

	case stepPort:
		sb.WriteString(stepNumber.Render("[Step 6/8] ") + wizardTitle.Render("Server Port"))
		sb.WriteString("\nHTTP and WebSocket API server port (default 3000):\n\n")
		sb.WriteString(promptStyle.Render("Port: ") + m.textInput.View() + "\n")
		if m.errMessage != "" {
			sb.WriteString("\n" + badgeDisconnected.Render("! "+m.errMessage) + "\n")
		}
		sb.WriteString("\n" + helpDescStyle.Render("Enter: continue | Esc: cancel"))

	case stepSkipOld:
		sb.WriteString(stepNumber.Render("[Step 7/8] ") + wizardTitle.Render("Skip Old Backlog Messages") + "\n")
		sb.WriteString("Skip messages received while the bot was offline?\n\n")

		options := []string{
			"Yes (Recommended - ignore old unread messages)",
			"No (Process all backlog messages)",
		}
		for i, opt := range options {
			radio := "( ) "
			text := unselectedItemStyle.Render(opt)
			if i == m.choiceIndex {
				radio = cursorStyle.Render("(*) ")
				text = selectedItemStyle.Render(opt)
			}
			sb.WriteString(fmt.Sprintf("  %s%s\n", radio, text))
		}
		sb.WriteString("\n" + helpDescStyle.Render("Up/Down: navigate | Enter: select | Esc: cancel"))

	case stepVerbose:
		sb.WriteString(stepNumber.Render("[Step 8/8] ") + wizardTitle.Render("Verbose Debug Logging") + "\n")
		sb.WriteString("Enable detailed debug logs?\n\n")

		options := []string{
			"Yes (Verbose debug logs)",
			"No (Clean standard logs - Recommended)",
		}
		for i, opt := range options {
			radio := "( ) "
			text := unselectedItemStyle.Render(opt)
			if i == m.choiceIndex {
				radio = cursorStyle.Render("(*) ")
				text = selectedItemStyle.Render(opt)
			}
			sb.WriteString(fmt.Sprintf("  %s%s\n", radio, text))
		}
		sb.WriteString("\n" + helpDescStyle.Render("Up/Down: navigate | Enter: select | Esc: cancel"))

	case stepSaveEnv:
		sb.WriteString(wizardTitle.Render("Save Configuration to .env") + "\n")
		sb.WriteString("Save this configuration to .env so WhatsRook starts automatically next time?\n\n")

		options := []string{
			"Yes (Save to .env)",
			"No (Run for this session only)",
		}
		for i, opt := range options {
			radio := "( ) "
			text := unselectedItemStyle.Render(opt)
			if i == m.choiceIndex {
				radio = cursorStyle.Render("(*) ")
				text = selectedItemStyle.Render(opt)
			}
			sb.WriteString(fmt.Sprintf("  %s%s\n", radio, text))
		}
		sb.WriteString("\n" + helpDescStyle.Render("Up/Down: navigate | Enter: select | Esc: cancel"))

	case stepSummary:
		sb.WriteString(wizardTitle.Render("Configuration Ready!") + "\n")

		authMethod := "Pair Code"
		if m.args.QRCode {
			authMethod = "QR Code"
		}
		cacheName := "In-Memory"
		if m.args.RedisURL != "" {
			cacheName = m.args.RedisURL
		}
		skipStr := "Yes"
		if !m.args.SkipOldMessages {
			skipStr = "No"
		}
		verbStr := "No"
		if m.args.Verbose {
			verbStr = "Yes"
		}

		sessionStr := m.args.Session
		if sessionStr == "" {
			sessionStr = "Standby (Idle Mode)"
		}

		summaryContent := fmt.Sprintf(
			"Session:           %s\nAuth Method:       %s\nClient Identity:   %s\nDatabase:          %s\nCache Backend:     %s\nServer Port:       %d\nSkip Old Msgs:     %s\nVerbose Logs:      %s",
			badgeConnected.Render(sessionStr),
			pillStyle.Render(authMethod),
			pillStyle.Render(m.args.Client),
			pillStyle.Render(m.args.Database),
			pillStyle.Render(cacheName),
			m.args.Port,
			pillStyle.Render(skipStr),
			pillStyle.Render(verbStr),
		)

		sb.WriteString(infoBoxStyle.Render(summaryContent))
		sb.WriteString("\n\n")

		cancelBtn := btnInactiveCancel.Render("Cancel")
		if m.summaryChoice == 0 {
			cancelBtn = btnActiveCancel.Render("Cancel")
		}

		confirmBtn := btnInactiveConfirm.Render("Confirm & Start")
		if m.summaryChoice == 1 {
			confirmBtn = btnActiveConfirm.Render("Confirm & Start")
		}

		actionButtons := lipgloss.JoinHorizontal(lipgloss.Center, cancelBtn, "    ", confirmBtn)
		sb.WriteString(actionButtons)
		sb.WriteString("\n\n")
		sb.WriteString(helpDescStyle.Render("Left/Right or Tab: choose | Enter: execute | Esc: cancel\n"))
	}

	cardWidth := 64
	if m.width > 0 && m.width < 68 {
		cardWidth = m.width - 4
	}
	card := modalBoxStyle.Width(cardWidth).Render(sb.String())
	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, card)
	}
	return card
}

// RunWizard launches the Bubbletea interactive setup wizard in alt screen mode.
func RunWizard(initial Config) (Config, bool) {
	model := NewWizardModel(initial)
	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return initial, false
	}
	if wm, ok := finalModel.(WizardModel); ok {
		return wm.args, wm.proceed
	}
	return initial, false
}

// CleanAndValidatePhone strips formatting characters and verifies digit count.
func CleanAndValidatePhone(input string) (string, bool) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return "", false
	}

	var digits strings.Builder
	hasPlus := strings.HasPrefix(trimmed, "+")

	for _, r := range trimmed {
		if unicode.IsDigit(r) {
			digits.WriteRune(r)
		}
	}

	numStr := digits.String()
	if len(numStr) < 7 || len(numStr) > 15 {
		return "", false
	}

	if hasPlus {
		return "+" + numStr, true
	}
	return numStr, true
}

// SaveEnvConfig persists selected configuration keys into a .env file.
func SaveEnvConfig(cfg Config, filename string) error {
	var envMap = make(map[string]string)

	if data, err := os.ReadFile(filename); err == nil {
		lines := strings.SplitSeq(string(data), "\n")
		for l := range lines {
			l = strings.TrimSpace(l)
			if l == "" || strings.HasPrefix(l, "#") {
				continue
			}
			if idx := strings.Index(l, "="); idx > 0 {
				k := strings.TrimSpace(l[:idx])
				v := strings.TrimSpace(l[idx+1:])
				v = strings.Trim(v, `"'`)
				envMap[k] = v
			}
		}
	}

	if cfg.Session != "" {
		envMap["SESSION"] = cfg.Session
	}
	if cfg.Client != "" {
		envMap["CLIENT"] = cfg.Client
	}
	if cfg.Database != "" {
		envMap["DATABASE_URL"] = cfg.Database
	}
	if cfg.RedisURL != "" {
		envMap["REDIS_URL"] = cfg.RedisURL
	}
	if cfg.Port > 0 {
		envMap["PORT"] = strconv.Itoa(cfg.Port)
	}
	if cfg.SkipOldMessages {
		envMap["SKIP_OLD_MESSAGES"] = "true"
	}
	if cfg.Verbose {
		envMap["VERBOSE"] = "true"
	}

	var out strings.Builder
	out.WriteString("# WhatsRook Configuration Generated by Setup Wizard\n")
	for k, v := range envMap {
		out.WriteString(fmt.Sprintf("%s=\"%s\"\n", k, v))
	}

	return os.WriteFile(filename, []byte(out.String()), 0600)
}
