package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/yd4dev/gubi/internal/bot"
	"github.com/yd4dev/gubi/internal/database"
	"github.com/yd4dev/gubi/pkg/platform/otakugifs"
)

func init() {
	bot.Register(&KissCommand{})
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

	imageURL, err := otakugifs.FetchGIF(otakugifs.ReactionKiss)
	if err != nil {
		return err
	}

	kisser := event.Member()
	kissed, _ := event.SlashCommandInteractionData().OptMember("member")

	kisses := database.Kisses{FirstID: min(kisser.User.ID, kissed.User.ID), SecondID: max(kisser.User.ID, kissed.User.ID)}

	database.DB.FirstOrCreate(&kisses)

	kisses.Kisses += 1

	database.DB.Save(&kisses)

	if err = event.CreateMessage(discord.NewMessageCreate().
		AddEmbeds(discord.NewEmbed().
			WithImage(imageURL).
			WithDescriptionf("%s just kissed %s!", kisser.Mention(), kissed.Mention()).
			WithFooterTextf("Total kisses: %d", kisses.Kisses))); err != nil {
		return err
	}
	return nil
}
