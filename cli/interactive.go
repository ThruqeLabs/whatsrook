package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// cleanAndValidatePhone strips formatting characters and verifies digit count.
func cleanAndValidatePhone(input string) (string, bool) {
	s := strings.TrimSpace(input)
	if s == "" {
		return "", false
	}
	hasPlus := strings.HasPrefix(s, "+")
	digits := make([]rune, 0, len(s))
	for _, r := range s {
		if r >= '0' && r <= '9' {
			digits = append(digits, r)
		}
	}
	if len(digits) < 7 || len(digits) > 15 {
		return "", false
	}
	result := string(digits)
	if hasPlus {
		result = "+" + result
	}
	return result, true
}

// readLine reads a single line from the buffered reader.
func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "", err
	}
	return strings.TrimRight(line, "\r\n"), nil
}

// RunInteractiveWizard prompts the user step-by-step to configure WhatsRook.
func RunInteractiveWizard(initial CLIArgs, in io.Reader, out io.Writer) (CLIArgs, bool) {
	reader := bufio.NewReader(in)
	res := initial

	_, _ = fmt.Fprintln(out, `
===============================================================
               WhatsRook Setup Wizard
   Welcome! Let's get your WhatsApp bot set up in seconds.
===============================================================`)

	// ── Step 1: Session / Phone Number
	for {
		_, _ = fmt.Fprintln(out, `
[1/7] WhatsApp Phone Number
---------------------------------------------------------------
Enter your WhatsApp phone number with country code
(e.g., +1234567890 or 2348012345678).

Or press Enter without typing anything to start in Standby (Idle) Mode.`)

		defaultPhone := initial.Session
		if defaultPhone != "" {
			_, _ = fmt.Fprintf(out, "Phone Number [%s]: ", defaultPhone)
		} else {
			_, _ = fmt.Fprint(out, "Phone Number [or Enter for Standby]: ")
		}

		line, err := readLine(reader)
		if err != nil {
			_, _ = fmt.Fprintln(out, "\nSetup cancelled.")
			return initial, false
		}
		line = strings.TrimSpace(line)

		if line == "" {
			if defaultPhone != "" {
				res.Session = defaultPhone
				break
			}
			// Ask if user wants Standby Mode
			_, _ = fmt.Fprint(out, "Start WhatsRook in Standby (Idle) Mode? [Y/n]: ")
			confLine, err := readLine(reader)
			if err != nil {
				return initial, false
			}
			confLine = strings.ToLower(strings.TrimSpace(confLine))
			if confLine == "" || confLine == "y" || confLine == "yes" {
				res.Session = ""
				break
			}
			continue
		}

		cleaned, ok := cleanAndValidatePhone(line)
		if !ok {
			_, _ = fmt.Fprintln(out, "(!) Invalid phone number. Please include country code (7-15 digits total).")
			continue
		}
		res.Session = cleaned
		break
	}

	// If Session is empty, we only need port before running idle mode
	if res.Session == "" {
		portDefault := initial.Port
		if portDefault <= 0 {
			portDefault = 3000
		}
		_, _ = fmt.Fprintf(out, `
[2/2] Standby HTTP Server Port
---------------------------------------------------------------
Port for the standby healthcheck server [default: %d]: `, portDefault)
		portLine, _ := readLine(reader)
		portLine = strings.TrimSpace(portLine)
		if portLine != "" {
			if p, err := strconv.Atoi(portLine); err == nil && p > 0 && p <= 65535 {
				res.Port = p
			}
		} else {
			res.Port = portDefault
		}

		_, _ = fmt.Fprintf(out, `
===============================================================
                   Configuration Summary
===============================================================
  Mode:       Standby (Idle Server)
  HTTP Port:  %d
===============================================================

Starting WhatsRook in Standby Mode...
`, res.Port)
		return res, true
	}

	// ── Step 2: Authentication / Pairing Method
	for {
		_, _ = fmt.Fprintln(out, `
[2/7] Pairing / Authentication Method
---------------------------------------------------------------
How would you like to link your WhatsApp account?

  1) Pair Code (Recommended)
     An 8-character pairing code will be displayed in the console.
     Open WhatsApp -> Linked Devices -> Link with phone number.

  2) QR Code
     A QR code will be generated and served via local web server
     for you to scan with WhatsApp camera.`)

		defaultChoice := "1"
		if initial.QRCode {
			defaultChoice = "2"
		}
		_, _ = fmt.Fprintf(out, "Select option [1-2] (default: %s): ", defaultChoice)

		line, err := readLine(reader)
		if err != nil {
			return initial, false
		}
		line = strings.TrimSpace(line)
		if line == "" {
			line = defaultChoice
		}

		if line == "1" || strings.EqualFold(line, "pair") {
			res.Pair = true
			res.QRCode = false
			break
		} else if line == "2" || strings.EqualFold(line, "qr") || strings.EqualFold(line, "qrcode") {
			res.Pair = false
			res.QRCode = true
			break
		} else {
			_, _ = fmt.Fprintln(out, "(!) Please select 1 (Pair Code) or 2 (QR Code).")
		}
	}

	// ── Step 3: Client Identity
	for {
		_, _ = fmt.Fprintln(out, `
[3/7] Client Device Identity
---------------------------------------------------------------
Select the platform type that WhatsRook will emulate:

  1) Chrome / Web (Default - Recommended)
  2) Android (Emulate Android phone)
  3) iOS (Emulate iPhone)`)

		defaultClientChoice := "1"
		switch strings.ToLower(initial.Client) {
		case "android":
			defaultClientChoice = "2"
		case "ios":
			defaultClientChoice = "3"
		default:
			defaultClientChoice = "1"
		}
		_, _ = fmt.Fprintf(out, "Select option [1-3] (default: %s): ", defaultClientChoice)

		line, err := readLine(reader)
		if err != nil {
			return initial, false
		}
		line = strings.TrimSpace(line)
		if line == "" {
			line = defaultClientChoice
		}

		if line == "1" || strings.EqualFold(line, "chrome") || strings.EqualFold(line, "web") {
			res.Client = "chrome"
			break
		} else if line == "2" || strings.EqualFold(line, "android") {
			res.Client = "android"
			break
		} else if line == "3" || strings.EqualFold(line, "ios") {
			res.Client = "ios"
			break
		} else {
			_, _ = fmt.Fprintln(out, "(!) Please select 1 (Chrome), 2 (Android), or 3 (iOS).")
		}
	}

	// ── Step 4: Database Storage
	_, _ = fmt.Fprintln(out, `
[4/7] Database Storage
---------------------------------------------------------------
WhatsRook stores session data and plugin settings.
By default, it uses SQLite (local 'whatsrook.db' file - zero setup).

If you want to use PostgreSQL, enter the connection URL.
(e.g., postgres://user:password@localhost:5432/whatsrook?sslmode=disable)
Otherwise, press Enter to use SQLite.`)

	defaultDB := initial.Database
	if defaultDB == "" {
		defaultDB = "sqlite"
	}
	_, _ = fmt.Fprintf(out, "Database URL [default: %s]: ", defaultDB)

	dbLine, err := readLine(reader)
	if err != nil {
		return initial, false
	}
	dbLine = strings.TrimSpace(dbLine)
	if dbLine == "" {
		res.Database = defaultDB
	} else {
		res.Database = dbLine
	}

	// ── Step 5: Redis Cache
	_, _ = fmt.Fprintln(out, `
[5/7] Cache Backend (Optional)
---------------------------------------------------------------
WhatsRook uses an in-memory cache by default.
If you have a Redis server, enter the Redis connection URL.
(e.g., redis://localhost:6379/0)
Otherwise, press Enter to use fast built-in in-memory cache.`)

	defaultRedis := initial.RedisURL
	redisPrompt := "Redis URL [press Enter to skip]: "
	if defaultRedis != "" {
		redisPrompt = fmt.Sprintf("Redis URL [%s]: ", defaultRedis)
	}
	_, _ = fmt.Fprint(out, redisPrompt)

	redisLine, err := readLine(reader)
	if err != nil {
		return initial, false
	}
	redisLine = strings.TrimSpace(redisLine)
	if redisLine == "" {
		res.RedisURL = defaultRedis
	} else if strings.EqualFold(redisLine, "skip") || strings.EqualFold(redisLine, "none") {
		res.RedisURL = ""
	} else {
		res.RedisURL = redisLine
	}

	// ── Step 6: Server Port
	portDefault := initial.Port
	if portDefault <= 0 {
		portDefault = 3000
	}
	for {
		_, _ = fmt.Fprintf(out, `
[6/7] WebSocket & Web Server Port
---------------------------------------------------------------
HTTP and WebSocket API server port [default: %d]: `, portDefault)

		portLine, err := readLine(reader)
		if err != nil {
			return initial, false
		}
		portLine = strings.TrimSpace(portLine)
		if portLine == "" {
			res.Port = portDefault
			break
		}
		if p, err := strconv.Atoi(portLine); err == nil && p > 0 && p <= 65535 {
			res.Port = p
			break
		}
		_, _ = fmt.Fprintln(out, "(!) Please enter a valid port number between 1 and 65535.")
	}

	// ── Step 7: Skip Old Messages & Verbose
	_, _ = fmt.Fprintln(out, `
[7/7] Additional Options
---------------------------------------------------------------`)

	// Skip Old Messages
	_, _ = fmt.Fprint(out, "* Skip messages sent while bot was offline? [Y/n] (default: Yes): ")
	skipLine, err := readLine(reader)
	if err != nil {
		return initial, false
	}
	skipLine = strings.ToLower(strings.TrimSpace(skipLine))
	if skipLine == "n" || skipLine == "no" || skipLine == "false" {
		res.SkipOldMessages = false
	} else {
		res.SkipOldMessages = true
	}

	// Verbose Logging
	_, _ = fmt.Fprint(out, "* Enable verbose debug logging? [y/N] (default: No): ")
	verbLine, err := readLine(reader)
	if err != nil {
		return initial, false
	}
	verbLine = strings.ToLower(strings.TrimSpace(verbLine))
	if verbLine == "y" || verbLine == "yes" || verbLine == "true" {
		res.Verbose = true
	} else {
		res.Verbose = false
	}

	// ── Step 8: Save to .env
	_, _ = fmt.Fprintln(out, `
---------------------------------------------------------------
Save this configuration to '.env' so WhatsRook starts
automatically with these settings in the future?`)
	_, _ = fmt.Fprint(out, "Save to .env? [Y/n] (default: Yes): ")
	saveLine, err := readLine(reader)
	if err != nil {
		return initial, false
	}
	saveLine = strings.ToLower(strings.TrimSpace(saveLine))
	if saveLine == "" || saveLine == "y" || saveLine == "yes" {
		if err := SaveEnvConfig(res, ".env"); err != nil {
			_, _ = fmt.Fprintf(out, "(!) Could not save .env: %v\n", err)
		} else {
			_, _ = fmt.Fprintln(out, "✓ Configuration successfully saved to .env")
		}
	}

	// ── Summary Card
	authMethod := "Pair Code"
	if res.QRCode {
		authMethod = "QR Code"
	}
	cacheBackend := "In-Memory"
	if res.RedisURL != "" {
		cacheBackend = res.RedisURL
	}
	skipOldStr := "Yes"
	if !res.SkipOldMessages {
		skipOldStr = "No"
	}
	verboseStr := "No"
	if res.Verbose {
		verboseStr = "Yes"
	}

	_, _ = fmt.Fprintf(out, `
===============================================================
                   Configuration Summary
===============================================================
  Phone / Session:   %s
  Auth Method:       %s
  Client Identity:   %s
  Database:          %s
  Cache Backend:     %s
  Server Port:       %d
  Skip Old Messages: %s
  Verbose Logs:      %s
===============================================================

Starting WhatsRook...
`, res.Session, authMethod, res.Client, res.Database, cacheBackend, res.Port, skipOldStr, verboseStr)

	return res, true
}

// SaveEnvConfig writes or updates .env with the given CLIArgs.
func SaveEnvConfig(args CLIArgs, filename string) error {
	if filename == "" {
		filename = ".env"
	}

	// Read existing .env if present
	existing := make(map[string]string)
	lines := make([]string, 0)

	if data, err := os.ReadFile(filename); err == nil {
		rawLines := strings.SplitSeq(string(data), "\n")
		for l := range rawLines {
			trimmed := strings.TrimSpace(l)
			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				lines = append(lines, l)
				continue
			}
			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) == 2 {
				k := strings.TrimSpace(parts[0])
				existing[k] = strings.TrimSpace(parts[1])
			}
			lines = append(lines, l)
		}
	}

	// Prepare updated values
	newVals := make(map[string]string)
	if args.Session != "" {
		newVals["SESSION"] = args.Session
	}
	if args.Pair {
		newVals["PAIR"] = "true"
		newVals["QRCODE"] = "false"
	} else if args.QRCode {
		newVals["PAIR"] = "false"
		newVals["QRCODE"] = "true"
	}
	if args.Client != "" {
		newVals["CLIENT"] = args.Client
	}
	if args.Database != "" {
		newVals["DATABASE_URL"] = args.Database
	}
	if args.RedisURL != "" {
		newVals["REDIS_URL"] = args.RedisURL
	}
	if args.Port > 0 {
		newVals["PORT"] = strconv.Itoa(args.Port)
	}
	if !args.SkipOldMessages {
		newVals["SKIP_OLD_MESSAGES"] = "false"
	} else {
		newVals["SKIP_OLD_MESSAGES"] = "true"
	}
	if args.Verbose {
		newVals["VERBOSE"] = "true"
	} else {
		newVals["VERBOSE"] = "false"
	}

	// Build output content
	var sb strings.Builder
	writtenKeys := make(map[string]bool)

	if len(lines) > 0 {
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				sb.WriteString(line)
				sb.WriteString("\n")
				continue
			}
			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) == 2 {
				k := strings.TrimSpace(parts[0])
				if v, ok := newVals[k]; ok {
					fmt.Fprintf(&sb, "%s=%s\n", k, v)
					writtenKeys[k] = true
					continue
				}
			}
			sb.WriteString(line)
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString("# WhatsRook Environment Configuration\n\n")
	}

	// Append any remaining keys not present in original file
	orderedKeys := []string{
		"SESSION", "PAIR", "QRCODE", "CLIENT", "DATABASE_URL",
		"REDIS_URL", "PORT", "SKIP_OLD_MESSAGES", "VERBOSE",
	}
	for _, k := range orderedKeys {
		if v, ok := newVals[k]; ok && !writtenKeys[k] {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}

	dir := filepath.Dir(filename)
	if dir != "" && dir != "." {
		_ = os.MkdirAll(dir, 0755)
	}
	return os.WriteFile(filename, []byte(sb.String()), 0644)
}
