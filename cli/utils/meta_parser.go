package cliutils

import (
	"whatsrook/cli/utils/meta"
)

// CommandInfo mirrors registered command information for Meta AI context.
type CommandInfo = meta.CommandInfo

var (
	// BuildRunCommandInstruction builds the instruction block prepended to every Meta AI request.
	BuildRunCommandInstruction = meta.BuildRunCommandInstruction

	// BuildRunCommandInstructionWithName builds the instruction block using custom botName.
	BuildRunCommandInstructionWithName = meta.BuildRunCommandInstructionWithName

	// BuildRunCommandInstructionWithNameAndPrefix builds the instruction block using custom botName and prefix.
	BuildRunCommandInstructionWithNameAndPrefix = meta.BuildRunCommandInstructionWithNameAndPrefix

	// ParseRunCommand checks whether an AI reply is requesting that the bot run a command.
	ParseRunCommand = meta.ParseRunCommand

	// AnswerParserString converts an AI-generated response written in Markdown into plain text.
	AnswerParserString = meta.AnswerParserString

	// RenderGroupContext turns GroupInfo into a text block appended to the Meta AI prompt.
	RenderGroupContext = meta.RenderGroupContext

	// RenderUserContext turns user info into a text block appended to the Meta AI prompt.
	RenderUserContext = meta.RenderUserContext

	// RenderQuotedContext turns quoted-message info into a text block for Meta AI context.
	RenderQuotedContext = meta.RenderQuotedContext
)
