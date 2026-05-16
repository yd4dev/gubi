package commands

import (
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/yd4dev/gubi/internal/bot"
	"github.com/yd4dev/gubi/pkg/platform/nekosbest"
)

func init() {
	bot.RegisterCommand(&SlapCommand{})
}

type SlapCommand struct{}

func (*SlapCommand) Definition() discord.ApplicationCommandCreate {
	return discord.SlashCommandCreate{
		Name:        "slap",
		Description: "Slap another user!",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionUser{
				Name:        "user",
				Description: "The user you want to slap.",
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

func (*SlapCommand) Handler(event *events.ApplicationCommandInteractionCreate) error {
	slapper := event.User()
	slapped, _ := event.SlashCommandInteractionData().OptUser("user")

	if slapper.ID == slapped.ID {
		return event.CreateMessage(discord.NewMessageCreate().WithContent("You cannot slap yourself! 😔").WithEphemeral(true))
	}

	gif, err := nekosbest.FetchGIF(nekosbest.Slap)
	if err != nil {
		return err
	}
	imageURL := gif.URL

	media := discord.MediaGalleryItem{
		Media: discord.UnfurledMediaItem{
			URL:         imageURL,
			ContentType: "image/gif",
		},
	}

	return event.CreateMessage(discord.NewMessageCreate().AddComponents(
		discord.NewContainer(
			discord.NewTextDisplayf("### %s just slapped %s!", discord.UserMention(slapper.ID), discord.UserMention(slapped.ID)),
			discord.NewMediaGallery(media),
			discord.NewTextDisplayf("Anime: %s", gif.AnimeName),
		).WithAccentColor(0x880808),
		discord.NewActionRow(
			discord.NewSecondaryButton("Search Anime", "anime:search:"+strings.ReplaceAll(gif.AnimeName, ":", " ")),
		),
	).WithFlags(discord.MessageFlagIsComponentsV2))
}
