package bot

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var RegisteredCommands map[string]Command = make(map[string]Command)
var DiscordApplicationCommandCreates []discord.ApplicationCommandCreate

type Command interface {
	Definition() discord.ApplicationCommandCreate
	Handler(event *events.ApplicationCommandInteractionCreate) error
}

func RegisterCommand(command Command) {
	RegisteredCommands[command.Definition().CommandName()] = command
	DiscordApplicationCommandCreates = append(DiscordApplicationCommandCreates, command.Definition())
}

func CommandHandler(event *events.ApplicationCommandInteractionCreate) {
	cmd := RegisteredCommands[event.Data.CommandName()]
	if cmd == nil {
		slog.Error("Command was not found.", slog.String("commandName", event.Data.CommandName()))
	}
	if err := cmd.Handler(event); err != nil {
		slog.Error("Error occured while running command.", slog.String("commandName", cmd.Definition().CommandName()), slog.Any("err", err))
		event.CreateMessage(discord.NewMessageCreate().WithContent("An error occured.").WithEphemeral(true))
	}
}
