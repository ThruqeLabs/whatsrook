package tui

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"
	"whatsrook/logger"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go.uber.org/zap/zapcore"
)

// BotController defines the interface for managing and querying the WhatsApp bot.
type BotController interface {
	StartupTime() time.Time
	IsConnected() bool
	Connect() error
	PairPhone(ctx context.Context, phone string) (string, error)
	Logout(ctx context.Context) error
	ClearSessionDB(ctx context.Context, jid string)
}

// StartBotFunc defines the signature for spawning a bot instance.
type StartBotFunc func(ctx context.Context, cfg Config, logWriter io.Writer, pairCodeChan chan<- string) (BotController, error)

// LogBatchMsg contains a batch of log lines to reduce re-render overhead.
type LogBatchMsg []string

// PairCodeMsg is sent when a pair code is issued.
type PairCodeMsg string

// StatusUpdateMsg is sent periodically for stats & uptime.
type StatusUpdateMsg struct{}

// NotifTickMsg triggers the 1-second countdown for notifications.
type NotifTickMsg struct{}

type activeModal int

const (
	modalNone activeModal = iota
	modalSwitchSession
	modalQRCode
	modalHelp
	modalLogoutConfirm
)

// DashboardModel is the Bubbletea model for the live Agentic CLI dashboard.
type DashboardModel struct {
	bot            BotController
	args           Config
	viewport       viewport.Model
	logBuffer      []string
	pairCode       string
	connected      bool
	linking        bool
	verbose        bool
	uptime         time.Duration
	modal          activeModal
	switchInput    textinput.Model
	width          int
	height         int
	ctx            context.Context
	cancel         context.CancelFunc
	logChan        chan string
	pairCodeChan   chan string
	restartChan    chan Config
	startBot       StartBotFunc
	quitting       bool
	notification   string
	notifCountdown int
}

// TUILogWriter captures output written by Logger and sends it into Bubbletea.
type TUILogWriter struct {
	logChan chan<- string
}

// NewTUILogWriter creates a new TUILogWriter wrapping a string channel.
func NewTUILogWriter(logChan chan<- string) *TUILogWriter {
	return &TUILogWriter{logChan: logChan}
}

func (w *TUILogWriter) Write(p []byte) (n int, err error) {
	line := string(p)
	line = strings.TrimRight(line, "\r\n")
	if line != "" {
		select {
		case w.logChan <- line:
		default:
		}
	}
	return len(p), nil
}

// NewDashboardModel initializes the live dashboard.
func NewDashboardModel(ctx context.Context, cancel context.CancelFunc, bot BotController, args Config, logChan chan string, pairCodeChan chan string, restartChan chan Config, startBot StartBotFunc) DashboardModel {
	vp := viewport.New(80, 20)
	vp.SetContent("Initializing WhatsRook Dashboard...\n")

	ti := textinput.New()
	ti.Placeholder = "+1234567890"
	ti.Prompt = "> "
	ti.PromptStyle = cursorStyle
	ti.TextStyle = selectedItemStyle
	ti.PlaceholderStyle = helpDescStyle
	ti.Cursor.SetChar("|")
	ti.Cursor.Style = cursorStyle
	ti.CharLimit = 20
	ti.Width = 30

	return DashboardModel{
		bot:          bot,
		args:         args,
		viewport:     vp,
		logBuffer:    make([]string, 0, 1000),
		verbose:      args.Verbose,
		modal:        modalNone,
		switchInput:  ti,
		ctx:          ctx,
		cancel:       cancel,
		logChan:      logChan,
		pairCodeChan: pairCodeChan,
		restartChan:  restartChan,
		startBot:     startBot,
	}
}

func (m DashboardModel) Init() tea.Cmd {
	return tea.Batch(
		m.listenForLogs(),
		m.listenForPairCode(),
		m.tickUptime(),
	)
}

func (m DashboardModel) listenForLogs() tea.Cmd {
	return func() tea.Msg {
		line, ok := <-m.logChan
		if !ok {
			return nil
		}
		// Drain any available lines already waiting in the channel (batching)
		lines := []string{line}
	drainLoop:
		for range 50 {
			select {
			case extra, ok := <-m.logChan:
				if ok && extra != "" {
					lines = append(lines, extra)
				}
			default:
				break drainLoop
			}
		}
		return LogBatchMsg(lines)
	}
}

func (m DashboardModel) listenForPairCode() tea.Cmd {
	return func() tea.Msg {
		code, ok := <-m.pairCodeChan
		if !ok {
			return nil
		}
		return PairCodeMsg(code)
	}
}

func (m DashboardModel) tickUptime() tea.Cmd {
	return tea.Tick(time.Second, func(_ time.Time) tea.Msg {
		return StatusUpdateMsg{}
	})
}

func (m *DashboardModel) showNotification(msg string) tea.Cmd {
	m.notification = msg
	m.notifCountdown = 3
	return tea.Tick(time.Second, func(_ time.Time) tea.Msg {
		return NotifTickMsg{}
	})
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.recalculateViewport()
		return m, nil

	case LogBatchMsg:
		vpWidth := m.viewport.Width
		if vpWidth <= 0 {
			vpWidth = 80
		}

		for _, rawLine := range msg {
			// Sanitize and wrap line so terminal never auto-wraps and displaces layout
			wrapped := WrapLogLine(rawLine, vpWidth)
			m.logBuffer = append(m.logBuffer, wrapped...)

			// Detect connection status strings in log
			lower := strings.ToLower(rawLine)
			if strings.Contains(lower, "connected") {
				m.connected = true
				m.linking = false
			} else if strings.Contains(lower, "pair code:") || strings.Contains(lower, "pair code issued") {
				m.linking = true
			} else if strings.Contains(lower, "logged out") {
				m.connected = false
				m.linking = false
			}
		}

		// Cap buffer size
		if len(m.logBuffer) > 2000 {
			m.logBuffer = m.logBuffer[len(m.logBuffer)-2000:]
		}

		m.viewport.SetContent(strings.Join(m.logBuffer, "\n"))
		m.viewport.GotoBottom()
		cmds = append(cmds, m.listenForLogs())

	case PairCodeMsg:
		m.pairCode = string(msg)
		m.linking = true
		m.recalculateViewport()
		cmd := m.showNotification(fmt.Sprintf("New WhatsApp Pairing Code: %s", string(msg)))
		cmds = append(cmds, cmd, m.listenForPairCode())

	case StatusUpdateMsg:
		if m.bot != nil {
			if st := m.bot.StartupTime(); !st.IsZero() {
				m.uptime = time.Since(st)
			}
			m.connected = m.bot.IsConnected()
		}
		cmds = append(cmds, m.tickUptime())

	case NotifTickMsg:
		if m.notification != "" {
			m.notifCountdown--
			if m.notifCountdown <= 0 {
				m.notification = ""
			} else {
				cmds = append(cmds, tea.Tick(time.Second, func(_ time.Time) tea.Msg {
					return NotifTickMsg{}
				}))
			}
		}

	case tea.KeyMsg:
		// Dismiss notification on Esc / Enter if visible
		if m.notification != "" && m.modal == modalNone {
			switch msg.String() {
			case "esc", "enter", "space":
				m.notification = ""
				return m, nil
			}
		}

		// Modal handling
		if m.modal != modalNone {
			switch msg.String() {
			case "esc":
				m.modal = modalNone
				m.switchInput.Blur()
				return m, nil

			case "enter":
				if m.modal == modalSwitchSession {
					newPhone := strings.TrimSpace(m.switchInput.Value())
					if cleaned, ok := CleanAndValidatePhone(newPhone); ok {
						m.args.Session = cleaned
						m.args.Pair = true
						m.args.QRCode = false
						m.pairCode = ""
						cmd := m.showNotification(fmt.Sprintf("Switching session to %s...", cleaned))
						m.modal = modalNone
						m.switchInput.Blur()
						if m.restartChan != nil {
							go func(args Config) {
								m.restartChan <- args
							}(m.args)
						}
						return m, cmd
					}
					cmd := m.showNotification("Invalid phone number format.")
					return m, cmd
				} else if m.modal == modalLogoutConfirm {
					m.modal = modalNone
					if m.bot != nil {
						go func() {
							_ = m.bot.Logout(context.Background())
							m.bot.ClearSessionDB(context.Background(), "")
						}()
					}
					cmd := m.showNotification("Session logged out.")
					return m, cmd
				} else {
					m.modal = modalNone
					return m, nil
				}
			}

			if m.modal == modalSwitchSession {
				var cmd tea.Cmd
				m.switchInput, cmd = m.switchInput.Update(msg)
				return m, cmd
			}
			return m, nil
		}

		// Normal dashboard keybindings
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			if m.cancel != nil {
				m.cancel()
			}
			return m, tea.Quit

		case "p", "P":
			// Request new pair code
			if m.bot != nil && m.args.Session != "" {
				cmd := m.showNotification("Requesting WhatsApp pairing code...")
				go func() {
					code, err := m.bot.PairPhone(context.Background(), m.args.Session)
					if err == nil && code != "" {
						select {
						case m.pairCodeChan <- code:
						default:
						}
					}
				}()
				return m, cmd
			}

		case "s", "S":
			// Switch session modal
			m.modal = modalSwitchSession
			m.switchInput.SetValue("")
			m.switchInput.Focus()
			return m, textinput.Blink

		case "r", "R":
			// Reconnect
			if m.bot != nil {
				cmd := m.showNotification("Reconnecting WhatsApp socket...")
				go func() {
					_ = m.bot.Connect()
				}()
				return m, cmd
			}

		case "c", "C":
			// Clear logs
			m.logBuffer = m.logBuffer[:0]
			m.viewport.SetContent("")
			cmd := m.showNotification("Logs cleared.")
			return m, cmd

		case "v", "V":
			// Toggle verbose
			m.verbose = !m.verbose
			var notif string
			if m.verbose {
				Logger.SetLevel(zapcore.DebugLevel)
				notif = "Verbose debug logging ENABLED."
			} else {
				Logger.SetLevel(zapcore.InfoLevel)
				notif = "Verbose debug logging DISABLED."
			}
			cmd := m.showNotification(notif)
			return m, cmd

		case "l", "L":
			// Logout confirm modal
			m.modal = modalLogoutConfirm

		case "?", "h", "H":
			// Help modal
			m.modal = modalHelp

		default:
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *DashboardModel) recalculateViewport() {
	if m.width <= 0 || m.height <= 0 {
		return
	}

	headerHeight := 4
	if m.width < 75 {
		headerHeight = 6 // 2-line responsive header
	}
	if m.pairCode != "" {
		headerHeight += 7 // pair code banner
	}
	footerHeight := 2

	vpHeight := max(m.height-headerHeight-footerHeight-2, 3)
	vpWidth := max(m.width-4, 20)

	m.viewport.Width = vpWidth
	m.viewport.Height = vpHeight
	m.viewport.SetContent(strings.Join(m.logBuffer, "\n"))
	m.viewport.GotoBottom()
}

func (m DashboardModel) View() string {
	if m.quitting {
		return "\n  Shutting down WhatsRook...\n\n"
	}

	// 1. Floating Notification Modal in the Middle of Screen with 3s Countdown
	if m.notification != "" && m.modal == modalNone {
		notifContent := fmt.Sprintf(
			"%s\n\n%s\n\n%s %s",
			notifTitle.Render("NOTIFICATION"),
			selectedItemStyle.Render(m.notification),
			notifCountStyle.Render(fmt.Sprintf("[%ds]", m.notifCountdown)),
			helpDescStyle.Render("Closing automatically | Esc/Enter to dismiss"),
		)
		box := notifModalBox.Render(notifContent)
		if m.width > 0 && m.height > 0 {
			return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
		}
		return box
	}

	// 2. Active Modals (Switch session, Help, Logout confirm)
	if m.modal != modalNone {
		modalContent := ""
		switch m.modal {
		case modalSwitchSession:
			modalContent = fmt.Sprintf(
				"%s\n\nEnter new WhatsApp phone number (with country code):\n\n%s\n\n%s",
				wizardTitle.Render("Switch WhatsApp Session"),
				m.switchInput.View(),
				helpDescStyle.Render("Enter: Switch | Esc: Cancel"),
			)

		case modalLogoutConfirm:
			modalContent = fmt.Sprintf(
				"%s\n\nAre you sure you want to log out session %s?\nThis will disconnect WhatsApp and remove local auth keys.\n\n%s",
				badgeDisconnected.Render("Confirm Logout"),
				m.args.Session,
				helpKeyStyle.Render("Enter: Confirm Logout")+" | "+helpDescStyle.Render("Esc: Cancel"),
			)

		case modalHelp:
			modalContent = fmt.Sprintf(
				"%s\n\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n\n%s",
				wizardTitle.Render("WhatsRook Keybindings"),
				helpKeyStyle.Render("p")+" - Request new WhatsApp 8-character pairing code",
				helpKeyStyle.Render("s")+" - Switch or add a new WhatsApp session number",
				helpKeyStyle.Render("r")+" - Reconnect WhatsApp socket connection",
				helpKeyStyle.Render("c")+" - Clear logs viewport buffer",
				helpKeyStyle.Render("v")+" - Toggle verbose debug logging dynamically",
				helpKeyStyle.Render("l")+" - Log out current session from WhatsApp and shared DB",
				helpKeyStyle.Render("Up/Down, PgUp/PgDn")+" - Scroll through log history",
				helpKeyStyle.Render("q / Ctrl+C")+" - Cleanly shut down WhatsRook",
				helpDescStyle.Render("Press Esc or Enter to close"),
			)
		}

		if modalContent != "" {
			box := modalBoxStyle.Render(modalContent)
			if m.width > 0 && m.height > 0 {
				return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
			}
			return box
		}
	}

	var sb strings.Builder

	// 3. Responsive Header Bar
	statusBadge := badgeDisconnected.Render("[Offline]")
	if m.connected {
		statusBadge = badgeConnected.Render("[Connected]")
	} else if m.linking {
		statusBadge = badgeConnecting.Render("[Linking]")
	}

	sessionDisplay := m.args.Session
	if sessionDisplay == "" {
		sessionDisplay = "Standby (Idle Mode)"
	}

	uptimeStr := fmt.Sprintf("%02d:%02d:%02d",
		int(m.uptime.Hours()),
		int(m.uptime.Minutes())%60,
		int(m.uptime.Seconds())%60,
	)

	var headerText string
	if m.width < 75 {
		// Responsive 2-line compact header
		headerText = fmt.Sprintf(
			"%s  Status: %s  Uptime: %s\nSession: %s  Platform: %s  Port: %s",
			titleStyle.Render(" WhatsRook "),
			statusBadge,
			pillStyle.Render(uptimeStr),
			badgeConnected.Render(sessionDisplay),
			pillStyle.Render(m.args.Client),
			pillStyle.Render(fmt.Sprintf(":%d", m.args.Port)),
		)
	} else {
		// Wide single-line header
		headerText = fmt.Sprintf(
			"%s   Session: %s   Platform: %s   Port: %s   Status: %s   Uptime: %s",
			titleStyle.Render(" WhatsRook "),
			badgeConnected.Render(sessionDisplay),
			pillStyle.Render(m.args.Client),
			pillStyle.Render(fmt.Sprintf(":%d", m.args.Port)),
			statusBadge,
			pillStyle.Render(uptimeStr),
		)
	}

	sb.WriteString(headerBoxStyle.Render(headerText))
	sb.WriteString("\n")

	// 4. Pair Code Card (if active)
	if m.pairCode != "" {
		formattedCode := FormatPairCodeDisplay(m.pairCode)
		codeCard := fmt.Sprintf(
			"  WHATSAPP PAIRING CODE  \n\n     %s     \n\n  1. Open WhatsApp -> Settings -> Linked Devices\n  2. Tap 'Link with phone number instead'\n  3. Enter the 8-character code shown above",
			pairCodeNumberStyle.Render(formattedCode),
		)
		sb.WriteString(pairCodeBoxStyle.Render(codeCard))
		sb.WriteString("\n")
	}

	// 5. Main Boxed Logs Viewport
	logTitle := fmt.Sprintf(" Real-Time Logs & Activity (%d lines) ", len(m.logBuffer))
	sb.WriteString(logBoxStyle.Render(
		wizardTitle.Render(logTitle) + "\n" + m.viewport.View(),
	))
	sb.WriteString("\n")

	// 6. Interactive Hotkeys Footer
	footerKeys := fmt.Sprintf(
		"%s %s  %s %s  %s %s  %s %s  %s %s  %s %s  %s %s  %s %s",
		helpKeyStyle.Render("[p]"), helpDescStyle.Render("Pair"),
		helpKeyStyle.Render("[s]"), helpDescStyle.Render("Switch Session"),
		helpKeyStyle.Render("[r]"), helpDescStyle.Render("Reconnect"),
		helpKeyStyle.Render("[c]"), helpDescStyle.Render("Clear"),
		helpKeyStyle.Render("[v]"), helpDescStyle.Render("Verbose"),
		helpKeyStyle.Render("[l]"), helpDescStyle.Render("Logout"),
		helpKeyStyle.Render("[?]"), helpDescStyle.Render("Help"),
		helpKeyStyle.Render("[q]"), helpDescStyle.Render("Quit"),
	)
	sb.WriteString(footerKeys)

	return sb.String()
}

// WrapLogLine ensures long log lines don't exceed maxWidth and wrap cleanly without terminal displacement.
func WrapLogLine(line string, maxWidth int) []string {
	if maxWidth <= 10 {
		maxWidth = 80
	}
	// Strip carriage returns
	line = strings.ReplaceAll(line, "\r", "")

	// Fast path for short lines
	if len(line) <= maxWidth {
		return []string{line}
	}

	var result []string
	runes := []rune(line)
	for len(runes) > maxWidth {
		result = append(result, string(runes[:maxWidth]))
		runes = runes[maxWidth:]
	}
	if len(runes) > 0 {
		result = append(result, string(runes))
	}
	return result
}

// FormatPairCodeDisplay formats an 8-character WhatsApp pairing code with readable spacing.
func FormatPairCodeDisplay(code string) string {
	code = strings.TrimSpace(code)
	if len(code) == 8 {
		return fmt.Sprintf("[  %s %s %s %s - %s %s %s %s  ]",
			string(code[0]), string(code[1]), string(code[2]), string(code[3]),
			string(code[4]), string(code[5]), string(code[6]), string(code[7]),
		)
	}
	return "[  " + code + "  ]"
}

// RunDashboard starts the bot inside the rich Agentic Bubbletea dashboard.
func RunDashboard(parentCtx context.Context, initial Config, startBot StartBotFunc) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	logChan := make(chan string, 1000)
	pairCodeChan := make(chan string, 10)
	restartChan := make(chan Config, 1)

	var botCancel context.CancelFunc
	var currentBot BotController

	startBotWrapper := func(cfg Config) error {
		if botCancel != nil {
			botCancel()
		}

		botCtx, bCancel := context.WithCancel(ctx)
		botCancel = bCancel

		logWriter := NewTUILogWriter(logChan)
		bot, err := startBot(botCtx, cfg, logWriter, pairCodeChan)
		if err != nil {
			return err
		}
		currentBot = bot
		return nil
	}

	if err := startBotWrapper(initial); err != nil {
		return err
	}

	dashboardModel := NewDashboardModel(ctx, cancel, currentBot, initial, logChan, pairCodeChan, restartChan, startBot)
	p := tea.NewProgram(dashboardModel, tea.WithAltScreen())

	// Listen for session switch restart requests from TUI
	go func() {
		for {
			select {
			case newCfg, ok := <-restartChan:
				if !ok {
					return
				}
				_ = startBotWrapper(newCfg)
			case <-ctx.Done():
				return
			}
		}
	}()

	_, err := p.Run()
	if botCancel != nil {
		botCancel()
	}
	return err
}
