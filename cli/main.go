//go:generate go run github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest -platform-specific=true

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"whatsrook/logger"

	"whatsrook"
	"whatsrook/cli/tui"
	"whatsrook/cli/updater"
	"whatsrook/utils/cache"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "update" {
		ctx := context.Background()
		current := updater.GetStoredChannel()
		isBeta := current == "beta"
		up := updater.New(updater.Options{
			Out:     os.Stdout,
			Channel: current,
		})

		subcmd := "check"
		if len(os.Args) > 2 {
			subcmd = os.Args[2]
		}

		switch subcmd {
		case "check":
			_, err := up.Check(ctx)
			if err != nil {
				Logger.Error("update check failed", "err", err)
				os.Exit(1)
			}
			return

		case "upgrade", "apply", "now":
			res, err := up.Upgrade(ctx, isBeta)
			if err != nil {
				Logger.Error("upgrade failed", "err", err)
				os.Exit(1)
			}
			if res.Updated {
				fmt.Println("==> Restarting process with new binary...")
				err := updater.RestartProcess()
				Logger.Error("failed to restart process", "err", err)
				os.Exit(1)
			}
			return

		default:
			fmt.Fprintf(os.Stderr, "Unknown update subcommand %q. Usage: whatsrook update [check|upgrade]\n", subcmd)
			os.Exit(1)
		}
	}

	args := parseCLIArgs()

	defer func() {
		_ = cache.Close()
	}()

	defer PromptExit()

	if args.Update {
		ctx := context.Background()
		current := updater.GetStoredChannel()
		requested := args.UpdateChannel // "stable", "beta", or ""

		if requested != "" {
			if requested == current {
				fmt.Printf("==> Already on the %s channel.\n", current)
			} else {
				fmt.Printf("==> Switching from %s to %s channel...\n", current, requested)
				if err := updater.SetStoredChannel(requested); err != nil {
					Logger.Error("failed to set channel", "err", err)
					PromptExit()
					os.Exit(1)
				}
				current = requested
			}
		}

		up := updater.New(updater.Options{
			Out:     os.Stdout,
			Channel: current,
		})

		isBeta := current == "beta"
		res, err := up.Upgrade(ctx, isBeta)
		if err != nil {
			Logger.Error("update failed", "err", err)
			PromptExit()
			os.Exit(1)
		}
		if res.Updated {
			fmt.Println("==> Restarting process with new binary...")
			err := updater.RestartProcess()
			Logger.Error("failed to restart process", "err", err)
			PromptExit()
			os.Exit(1)
		}

		if args.Session == "" && os.Getenv("SESSION") == "" {
			fmt.Println("==> No active session requested. Exiting.")
			return
		}
	}

	// Interactive Wizard and Terminal Spawning Check
	if !args.Idle && !args.Logout && !args.Update {
		if args.Interactive || (args.Session == "" && IsInteractiveTerminal()) {
			SetInteractiveSession(true)
			tuiCfg := cliArgsToTUI(args)
			wizardCfg, proceed := tui.RunWizard(tuiCfg)
			if !proceed {
				return
			}
			args = tuiToCLIArgs(wizardCfg, args)
		} else if args.Session == "" && !IsInteractiveTerminal() && !IsCIOrDaemon() {
			spawned, err := SpawnTerminal()
			if err == nil && spawned {
				return
			}
		}
	}

	if args.RedisURL != "" {
		cache.Init(args.RedisURL)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Launch Agentic TUI Dashboard if in an interactive terminal and not disabled
	shouldRunTUI := !args.NoTUI && !args.Logout && (args.TUI || IsInteractiveTerminal() || IsInteractiveSession()) && !IsCIOrDaemon()

	if shouldRunTUI {
		startBotFunc := func(botCtx context.Context, cfg tui.Config, logWriter io.Writer, pairCodeChan chan<- string) (tui.BotController, error) {
			clientType, ok := whatsrook.ParseClientType(cfg.Client)
			if !ok {
				clientType = whatsrook.ClientChrome
			}

			bot := NewBot(BotConfig{
				Session:         cfg.Session,
				Pair:            cfg.Pair,
				QRCode:          cfg.QRCode,
				Logout:          cfg.Logout,
				Verbose:         cfg.Verbose,
				ClientType:      clientType,
				Database:        cfg.Database,
				WSPort:          cfg.Port,
				SkipOldMessages: cfg.SkipOldMessages,
				AsyncMessageAck: true,
				ConsoleOut:      logWriter,
			})

			go func() {
				if cfg.Session == "" {
					_ = runIdleMode(botCtx, cfg.Port)
				} else {
					_ = bot.Start(botCtx)
				}
			}()

			return bot, nil
		}

		if err := tui.RunDashboard(ctx, cliArgsToTUI(args), startBotFunc); err != nil && !errors.Is(err, context.Canceled) {
			Logger.Error("dashboard error", "err", err)
			PromptExit()
			os.Exit(1)
		}
		return
	}

	if args.Session == "" {
		if err := runIdleMode(ctx, args.Port); err != nil {
			Logger.Error("idle server error", "err", err)
			PromptExit()
			os.Exit(1)
		}
		return
	}

	clientType, ok := whatsrook.ParseClientType(args.Client)
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: unknown --client %q. Valid options: chrome, android, ios\n", args.Client)
		PromptExit()
		os.Exit(1)
	}

	bot := NewBot(BotConfig{
		Session:         args.Session,
		Pair:            args.Pair,
		QRCode:          args.QRCode,
		Logout:          args.Logout,
		Verbose:         args.Verbose,
		ClientType:      clientType,
		Database:        args.Database,
		WSPort:          args.Port,
		SkipOldMessages: args.SkipOldMessages,
		AsyncMessageAck: true,
	})

	if err := bot.Start(ctx); err != nil {
		if errors.Is(err, whatsrook.ErrLoggedOut) {
			if ctx.Err() != nil {
				return
			}
			Logger.Info("Session was logged out and removed. Switching to idle standby mode...")
			if err := runIdleMode(ctx, args.Port); err != nil {
				Logger.Error("idle server error", "err", err)
				PromptExit()
				os.Exit(1)
			}
			return
		}
		Logger.Error("bot error", "err", err)
		PromptExit()
		os.Exit(1)
	}
}

func runIdleMode(ctx context.Context, port int) error {
	if port <= 0 {
		port = 3000
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "WhatsRook standby • waiting for session • %s\n", time.Now().Format("15:04:05"))
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintln(w, "OK")
	})

	var listener net.Listener
	var actualPort int
	for p := port; p < port+100; p++ {
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", p))
		if err == nil {
			listener = l
			actualPort = p
			break
		}
		if p == port {
			Logger.Warn("port in use, attempting to bind alternative port", "attempted_port", p, "err", err)
		}
	}
	if listener == nil {
		return errors.New("failed to find an available port to bind HTTP server")
	}

	if actualPort != port {
		Logger.Warn("port in use — switched to alternative port", "original_port", port, "new_port", actualPort)
	}

	server := &http.Server{Handler: mux}
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			Logger.Error("http server error", "err", err)
		}
	}()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
		if listener != nil {
			_ = listener.Close()
		}
	}()

	fmt.Printf("\rWhatsRook standby • waiting for session • %s", time.Now().Format("15:04:05"))

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("\rWhatsRook standby • waiting for session • %s", time.Now().Format("15:04:05"))
		case <-ctx.Done():
			fmt.Println()
			return nil
		}
	}
}

func cliArgsToTUI(args CLIArgs) tui.Config {
	return tui.Config{
		Session:         args.Session,
		Pair:            args.Pair,
		QRCode:          args.QRCode,
		Client:          args.Client,
		Database:        args.Database,
		RedisURL:        args.RedisURL,
		Port:            args.Port,
		SkipOldMessages: args.SkipOldMessages,
		Verbose:         args.Verbose,
		Interactive:     args.Interactive,
		Idle:            args.Idle,
		Logout:          args.Logout,
		NoTUI:           args.NoTUI,
		TUI:             args.TUI,
	}
}

func tuiToCLIArgs(cfg tui.Config, original CLIArgs) CLIArgs {
	original.Session = cfg.Session
	original.Pair = cfg.Pair
	original.QRCode = cfg.QRCode
	original.Client = cfg.Client
	original.Database = cfg.Database
	original.RedisURL = cfg.RedisURL
	original.Port = cfg.Port
	original.SkipOldMessages = cfg.SkipOldMessages
	original.Verbose = cfg.Verbose
	original.Interactive = cfg.Interactive
	original.Idle = cfg.Idle
	original.Logout = cfg.Logout
	return original
}
