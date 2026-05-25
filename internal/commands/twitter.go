package commands

import (
	"regexp"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/yd4dev/gubi/internal/bot"
	"github.com/yd4dev/gubi/internal/shared"
)

func init() {
	bot.RegisterCommand(&TwitterCommand{})
}

type TwitterCommand struct{}

func (*TwitterCommand) Definition() discord.ApplicationCommandCreate {
	return discord.SlashCommandCreate{
		Name:        "twitter",
		Description: "Embed content from Twitter!",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "url",
				Description: "Twitter URL to fetch content from.",
				Required:    true,
			},
		},
		IntegrationTypes: []discord.ApplicationIntegrationType{
			//discord.ApplicationIntegrationTypeGuildInstall,
			discord.ApplicationIntegrationTypeUserInstall,
		},
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeBotDM,
			//discord.InteractionContextTypeGuild,
			discord.InteractionContextTypePrivateChannel,
		},
	}
}

func (*TwitterCommand) Handler(event *events.ApplicationCommandInteractionCreate) error {
	url, _ := event.SlashCommandInteractionData().OptString("url")
	r, err := regexp.Compile(`/status/(\d+)`)
	if err != nil {
		return err
	}

	event.DeferCreateMessage(true)

	if !r.MatchString(url) {
		_, err := event.Client().Rest.CreateFollowupMessage(event.Client().ApplicationID, event.Token(), bot.ErrorMessage("Invalid Twitter URL!"))
		return err
	}

	tweetID := r.FindStringSubmatch(url)[1]

	return shared.FetchAndDisplayTweet(event.Client(), event.Token(), tweetID)
}
