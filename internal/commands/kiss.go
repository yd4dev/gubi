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

func (_ *KissCommand) Definition() discord.ApplicationCommandCreate {
	return discord.SlashCommandCreate{
		Name:        "kiss",
		Description: "Kiss another member of this server!",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionUser{
				Name:        "member",
				Description: "The member you want to kiss.",
				Required:    true,
			},
		},
	}
}

func (_ *KissCommand) Handler(event *events.ApplicationCommandInteractionCreate) error {
	kisser := event.Member()
	kissed, _ := event.SlashCommandInteractionData().OptMember("member")

	message, err := shared.Kiss(kisser.User.ID, kissed.User.ID)
	if err != nil {
		return err
	}

	if err := event.CreateMessage(message); err != nil {
		return err
	}
	return nil
}
