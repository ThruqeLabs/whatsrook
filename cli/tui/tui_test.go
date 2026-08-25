package tui

import (
	"context"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type mockBotController struct {
	connected bool
	startTime time.Time
}

func (m *mockBotController) StartupTime() time.Time {
	return m.startTime
}

func (m *mockBotController) IsConnected() bool {
	return m.connected
}

func (m *mockBotController) Connect() error {
	m.connected = true
	return nil
}

func (m *mockBotController) PairPhone(_ context.Context, _ string) (string, error) {
	return "ABCD1234", nil
}

func (m *mockBotController) Logout(_ context.Context) error {
	m.connected = false
	return nil
}

func (m *mockBotController) ClearSessionDB(_ context.Context, _ string) {}

func TestWizardModel_ExistingEnvMenu(t *testing.T) {
	initial := Config{
		Session: "+1234567890",
		Client:  "chrome",
		Port:    3000,
	}

	model := NewWizardModel(initial)

	if model.step != stepExistingEnvMenu {
		t.Fatalf("expected stepExistingEnvMenu, got %v", model.step)
	}

	// 1. Test Quick Start (Enter on index 0)
	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := updated.(WizardModel)
	if !m.proceed {
		t.Errorf("expected proceed=true on Quick Start")
	}
	if cmd == nil {
		t.Errorf("expected quit cmd on Quick Start")
	}

	// 2. Test Navigate Down to Modify (index 1)
	model2 := NewWizardModel(initial)
	updated2, _ := model2.Update(tea.KeyMsg{Type: tea.KeyDown})
	m2 := updated2.(WizardModel)
	if m2.menuIndex != 1 {
		t.Errorf("expected menuIndex=1 after down arrow, got %d", m2.menuIndex)
	}

	updated3, _ := m2.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m3 := updated3.(WizardModel)
	if m3.step != stepPhone {
		t.Errorf("expected stepPhone after choosing Modify, got %v", m3.step)
	}
}

func TestWizardModel_StepProgression(t *testing.T) {
	initial := Config{
		Interactive: true,
	}

	model := NewWizardModel(initial)
	if model.step != stepPhone {
		t.Fatalf("expected start at stepPhone, got %v", model.step)
	}

	// Step 1: Enter phone
	model.textInput.SetValue("+1234567890")
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := updated.(WizardModel)
	if m.step != stepAuthMethod {
		t.Fatalf("expected stepAuthMethod, got %v", m.step)
	}
	if m.args.Session != "+1234567890" {
		t.Errorf("expected session +1234567890, got %q", m.args.Session)
	}

	// Step 2: Auth method (default choice 0 = Pair code)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(WizardModel)
	if m.step != stepClient {
		t.Fatalf("expected stepClient, got %v", m.step)
	}
	if !m.args.Pair {
		t.Errorf("expected Pair=true")
	}

	// Step 3: Client (select choice 1 = Android)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(WizardModel)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(WizardModel)
	if m.step != stepDatabase {
		t.Fatalf("expected stepDatabase, got %v", m.step)
	}
	if m.args.Client != "android" {
		t.Errorf("expected Client=android, got %q", m.args.Client)
	}

	// Step 4: Database (enter sqlite)
	m.textInput.SetValue("sqlite")
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(WizardModel)
	if m.step != stepRedis {
		t.Fatalf("expected stepRedis, got %v", m.step)
	}

	// Step 5: Redis (skip)
	m.textInput.SetValue("")
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(WizardModel)
	if m.step != stepPort {
		t.Fatalf("expected stepPort, got %v", m.step)
	}

	// Step 6: Port (3000)
	m.textInput.SetValue("3000")
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(WizardModel)
	if m.step != stepSkipOld {
		t.Fatalf("expected stepSkipOld, got %v", m.step)
	}

	// Step 7: Skip Old (choice 0 = Yes)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(WizardModel)
	if m.step != stepVerbose {
		t.Fatalf("expected stepVerbose, got %v", m.step)
	}

	// Step 8: Verbose (choice 1 = No)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(WizardModel)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(WizardModel)
	if m.step != stepSaveEnv {
		t.Fatalf("expected stepSaveEnv, got %v", m.step)
	}

	// Step 9: Save Env (choice 1 = No, avoid file write in unit test)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(WizardModel)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(WizardModel)
	if m.step != stepSummary {
		t.Fatalf("expected stepSummary, got %v", m.step)
	}

	// Step 10: Summary Navigation (Left -> Cancel, Right -> Confirm)
	if m.summaryChoice != 1 {
		t.Errorf("expected default summaryChoice=1 (Confirm), got %d", m.summaryChoice)
	}

	// Navigate Left (Cancel)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	m = updated.(WizardModel)
	if m.summaryChoice != 0 {
		t.Errorf("expected summaryChoice=0 after KeyLeft, got %d", m.summaryChoice)
	}

	// Toggle with Tab (Back to Confirm)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(WizardModel)
	if m.summaryChoice != 1 {
		t.Errorf("expected summaryChoice=1 after KeyTab, got %d", m.summaryChoice)
	}

	// Confirm & Start
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(WizardModel)
	if !m.proceed {
		t.Errorf("expected proceed=true after summary confirm")
	}
}

func TestWizardModel_PhoneValidation(t *testing.T) {
	initial := Config{Interactive: true}
	model := NewWizardModel(initial)

	// Invalid phone
	model.textInput.SetValue("123")
	updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := updated.(WizardModel)

	if m.step != stepPhone {
		t.Errorf("expected to stay at stepPhone on invalid phone")
	}
	if m.errMessage == "" {
		t.Errorf("expected errMessage to be set on invalid phone")
	}
}

func TestWizardModel_ViewRendering(t *testing.T) {
	model := NewWizardModel(Config{Interactive: true})
	view := model.View()
	if !strings.Contains(view, "WhatsRook") {
		t.Errorf("expected view to contain 'WhatsRook'")
	}
}

func TestFormatPairCodeDisplay(t *testing.T) {
	code := "8K3MW29P"
	formatted := FormatPairCodeDisplay(code)
	if !strings.Contains(formatted, "8 K 3 M - W 2 9 P") {
		t.Errorf("expected formatted pair code, got %q", formatted)
	}
}

func TestTUILogWriter(t *testing.T) {
	ch := make(chan string, 5)
	writer := NewTUILogWriter(ch)

	_, err := writer.Write([]byte("Test log line\n"))
	if err != nil {
		t.Fatalf("unexpected error on write: %v", err)
	}

	select {
	case line := <-ch:
		if line != "Test log line" {
			t.Errorf("expected 'Test log line', got %q", line)
		}
	default:
		t.Fatal("expected log message in channel")
	}
}

func TestDashboardModel_LogAndKeyHandling(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logChan := make(chan string, 10)
	pairCodeChan := make(chan string, 5)
	restartChan := make(chan Config, 1)

	args := Config{
		Session: "+1234567890",
		Client:  "chrome",
		Port:    3000,
	}

	mockBot := &mockBotController{connected: false, startTime: time.Now()}
	dashboard := NewDashboardModel(ctx, cancel, mockBot, args, logChan, pairCodeChan, restartChan, nil)

	// Test WindowSizeMsg (responsiveness)
	updated, _ := dashboard.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	d := updated.(DashboardModel)
	if d.width != 100 || d.height != 30 {
		t.Errorf("expected width=100, height=30, got width=%d, height=%d", d.width, d.height)
	}

	// Test LogBatchMsg update
	updated, _ = d.Update(LogBatchMsg([]string{"00:00:00 [INFO] Client connected"}))
	d = updated.(DashboardModel)
	if len(d.logBuffer) != 1 {
		t.Errorf("expected 1 log line in buffer, got %d", len(d.logBuffer))
	}
	if !d.connected {
		t.Errorf("expected connected=true after log detection")
	}

	// Test PairCodeMsg update
	updated, _ = d.Update(PairCodeMsg("ABCD1234"))
	d = updated.(DashboardModel)
	if d.pairCode != "ABCD1234" {
		t.Errorf("expected pairCode=ABCD1234, got %q", d.pairCode)
	}
	if d.notification == "" || d.notifCountdown != 3 {
		t.Errorf("expected notification countdown modal, got notif=%q, countdown=%d", d.notification, d.notifCountdown)
	}

	// Test 3s Notification Countdown Ticks
	updated, _ = d.Update(NotifTickMsg{})
	d = updated.(DashboardModel)
	if d.notifCountdown != 2 {
		t.Errorf("expected countdown=2, got %d", d.notifCountdown)
	}

	updated, _ = d.Update(NotifTickMsg{})
	d = updated.(DashboardModel)
	if d.notifCountdown != 1 {
		t.Errorf("expected countdown=1, got %d", d.notifCountdown)
	}

	updated, _ = d.Update(NotifTickMsg{})
	d = updated.(DashboardModel)
	if d.notification != "" {
		t.Errorf("expected notification to dismiss after 3 ticks, got %q", d.notification)
	}

	// Test key 'c' (clear logs and show notification)
	updated, _ = d.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	d = updated.(DashboardModel)
	if len(d.logBuffer) != 0 {
		t.Errorf("expected 0 log lines after clear, got %d", len(d.logBuffer))
	}
	if !strings.Contains(d.notification, "Logs cleared") {
		t.Errorf("expected notification for clear logs, got %q", d.notification)
	}

	// Dismiss notification via Esc
	updated, _ = d.Update(tea.KeyMsg{Type: tea.KeyEsc})
	d = updated.(DashboardModel)
	if d.notification != "" {
		t.Errorf("expected notification dismissed on Esc")
	}

	// Test key 'v' (toggle verbose)
	updated, _ = d.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})
	d = updated.(DashboardModel)
	if !d.verbose {
		t.Errorf("expected verbose=true after toggle")
	}

	// Dismiss notification via Enter
	updated, _ = d.Update(tea.KeyMsg{Type: tea.KeyEnter})
	d = updated.(DashboardModel)

	// Test key '?' (help modal)
	updated, _ = d.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	d = updated.(DashboardModel)
	if d.modal != modalHelp {
		t.Errorf("expected modalHelp, got %v", d.modal)
	}

	// Test Esc key (close modal)
	updated, _ = d.Update(tea.KeyMsg{Type: tea.KeyEsc})
	d = updated.(DashboardModel)
	if d.modal != modalNone {
		t.Errorf("expected modalNone after Esc, got %v", d.modal)
	}

	// Test View rendering in normal width
	view := d.View()
	if !strings.Contains(view, "WhatsRook") {
		t.Errorf("expected dashboard view to contain 'WhatsRook'")
	}
	if !strings.Contains(view, "+1234567890") {
		t.Errorf("expected dashboard view to contain session '+1234567890'")
	}

	// Test View rendering in narrow width (< 75)
	updated, _ = d.Update(tea.WindowSizeMsg{Width: 60, Height: 24})
	d = updated.(DashboardModel)
	narrowView := d.View()
	if !strings.Contains(narrowView, "WhatsRook") {
		t.Errorf("expected narrow dashboard view to contain 'WhatsRook'")
	}
}

func TestWrapLogLine(t *testing.T) {
	shortLine := "Simple log line"
	wrappedShort := WrapLogLine(shortLine, 80)
	if len(wrappedShort) != 1 || wrappedShort[0] != shortLine {
		t.Errorf("expected 1 unwrapped line, got %v", wrappedShort)
	}

	longLine := strings.Repeat("A", 150)
	wrappedLong := WrapLogLine(longLine, 60)
	if len(wrappedLong) != 3 {
		t.Errorf("expected 3 wrapped lines for length 150 at width 60, got %d", len(wrappedLong))
	}
	if len(wrappedLong[0]) != 60 || len(wrappedLong[1]) != 60 || len(wrappedLong[2]) != 30 {
		t.Errorf("expected chunk lengths [60, 60, 30], got lengths [%d, %d, %d]",
			len(wrappedLong[0]), len(wrappedLong[1]), len(wrappedLong[2]))
	}
}
