package shared

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/yd4dev/gubi/internal/database"
	"github.com/yd4dev/gubi/pkg/platform/nekosbest"
)

func Kiss(kisserID, kissedID snowflake.ID) (discord.MessageCreate, error) {
	if kisserID == kissedID {
		return discord.NewMessageCreate().WithContent("You cannot kiss yourself! 😛").WithEphemeral(true), nil
	}

	gif, err := nekosbest.FetchGIF(nekosbest.Kiss)
	if err != nil {
		return discord.MessageCreate{}, err
	}
	imageURL := gif.URL

	kisses := database.Kisses{FirstID: min(kisserID, kissedID), SecondID: max(kisserID, kissedID)}

	database.DB.FirstOrCreate(&kisses)

	kisses.Kisses += 1

	database.DB.Save(&kisses)

	media := discord.MediaGalleryItem{
		Media: discord.UnfurledMediaItem{
			URL:         imageURL,
			ContentType: "image/gif",
		},
	}

	return discord.NewMessageCreate().AddComponents(
		discord.NewContainer(
			discord.NewTextDisplayf("### %s just kissed %s!", discord.UserMention(kisserID), discord.UserMention(kissedID)),
			discord.NewMediaGallery(media),
			discord.NewTextDisplayf("Anime: %s | Total Kisses: %d", gif.AnimeName, kisses.Kisses),
		).WithAccentColor(0xFF6ECF),
		discord.NewActionRow(
			discord.NewPrimaryButton("Kiss back", "kiss_"+kissedID.String()+"_"+kisserID.String()).WithEmoji(discord.NewComponentEmoji("❤️")),
		),
	).WithFlags(discord.MessageFlagIsComponentsV2), nil
}
