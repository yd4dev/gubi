package bot

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var RegisteredCommands []Command
var DiscordApplicationCommandCreates []discord.ApplicationCommandCreate

type Command interface {
	Definition() discord.ApplicationCommandCreate
	Handler(event *events.ApplicationCommandInteractionCreate) error
}

func Register(command Command) {
	RegisteredCommands = append(RegisteredCommands, command)
	DiscordApplicationCommandCreates = append(DiscordApplicationCommandCreates, command.Definition())
}

func CommandHandler(event *events.ApplicationCommandInteractionCreate) {
	for _, cmd := range RegisteredCommands {
		if cmd.Definition().CommandName() == event.Data.CommandName() {
			if err := cmd.Handler(event); err != nil {
				slog.Error("Error occured while running command.", slog.String("commandName", cmd.Definition().CommandName()), slog.Any("err", err))
			}
			return
		}
	}
	slog.Error("Command was not found.", slog.String("commandName", event.Data.CommandName()))
}
