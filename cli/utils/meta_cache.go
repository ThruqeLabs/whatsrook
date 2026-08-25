package cliutils

import (
	"whatsrook/cli/utils/meta"
)

var (
	// GetOrBuildInstruction returns the cached instruction block if valid.
	GetOrBuildInstruction = meta.GetOrBuildInstruction

	// GetOrBuildInstructionWithName returns the cached instruction block for botName.
	GetOrBuildInstructionWithName = meta.GetOrBuildInstructionWithName

	// GetOrBuildInstructionWithNameAndPrefix returns the cached instruction block for botName and prefix.
	GetOrBuildInstructionWithNameAndPrefix = meta.GetOrBuildInstructionWithNameAndPrefix

	// ClearInstructionCache invalidates the cached RUN_COMMAND prompt block.
	ClearInstructionCache = meta.ClearInstructionCache

	// GetOrFetchGroupMeta returns cached GroupInfo for chatKey.
	GetOrFetchGroupMeta = meta.GetOrFetchGroupMeta
)
