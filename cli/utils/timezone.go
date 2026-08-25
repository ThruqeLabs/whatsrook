package cliutils

import (
	"whatsrook/cli/utils/timezone"
)

// WindowsTZEntry represents a Windows timezone definition mapped to IANA zones.
type WindowsTZEntry = timezone.WindowsTZEntry

var (
	// ResolveTimezoneAlias resolves user input to a canonical IANA timezone string.
	ResolveTimezoneAlias = timezone.ResolveTimezoneAlias
)
