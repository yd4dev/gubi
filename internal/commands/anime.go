package commands

import (
	"errors"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/yd4dev/gubi/internal/bot"
	"github.com/yd4dev/gubi/internal/shared"
)

func init() {
	bot.RegisterCommand(&AnimeCommand{})
}

type AnimeCommand struct{}

func (*AnimeCommand) Definition() discord.ApplicationCommandCreate {
	return discord.SlashCommandCreate{
		Name:        "anime",
		Description: "Find out a lot about anime!",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        "search",
				Description: "Search for an anime by name.",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:        "name",
						Description: "The name of the anime you want to search for.",
						Required:    true,
					},
				},
			},
		},
		IntegrationTypes: []discord.ApplicationIntegrationType{
			discord.ApplicationIntegrationTypeGuildInstall,
			discord.ApplicationIntegrationTypeUserInstall,
		},
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeBotDM,
			discord.InteractionContextTypeGuild,
			discord.InteractionContextTypePrivateChannel,
		},
	}
}

func (*AnimeCommand) Handler(event *events.ApplicationCommandInteractionCreate) error {
	subCommand := *event.SlashCommandInteractionData().SubCommandName

	switch subCommand {
	case "search":
		return search(event)
	}
	return errors.New("Subcommand " + subCommand + " not found.")
}

func search(event *events.ApplicationCommandInteractionCreate) error {
	name, _ := event.SlashCommandInteractionData().OptString("name")

	anime, err := shared.SearchAnimeByName(name)

	if err != nil {
		return err
	}

	if anime == nil {
		return event.CreateMessage(discord.NewMessageCreate().WithContentf("No anime found with name %s", name).WithEphemeral(true))
	}

	return event.CreateMessage(shared.DisplayAnime(anime))
}
