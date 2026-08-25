package cliutils

import (
	"whatsrook/cli/utils/meta"
)

// Data contains the full context for a Meta AI request.
type Data = meta.Data

// Tools describes which tools the AI may invoke in its response.
type Tools = meta.Tools

// Response is the structured reply from Meta AI.
type Response = meta.Response
