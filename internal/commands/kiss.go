package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/yd4dev/gubi/internal/bot"
	"github.com/yd4dev/gubi/internal/shared"
)

func init() {
	bot.RegisterCommand(&KissCommand{})
}

type KissCommand struct{}

func (*KissCommand) Definition() discord.ApplicationCommandCreate {
	return discord.SlashCommandCreate{
		Name:        "kiss",
		Description: "Kiss another user!",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionUser{
				Name:        "user",
				Description: "The user you want to kiss.",
				Required:    true,
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

func (*KissCommand) Handler(event *events.ApplicationCommandInteractionCreate) error {
	kisser := event.User()
	kissed, _ := event.SlashCommandInteractionData().OptUser("user")

	message, err := shared.Kiss(kisser.ID, kissed.ID)
	if err != nil {
		return err
	}

	if err := event.CreateMessage(message); err != nil {
		return err
	}
	return nil
}
