package shared

import (
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/rl404/verniy"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

var client = verniy.New()

func SearchAnimeByName(name string) (*verniy.Media, error) {
	res, err := client.SearchAnime(
		verniy.PageParamMedia{Search: name},
		1,
		1,
		verniy.MediaFieldBannerImage,
		verniy.MediaFieldTitle(verniy.MediaTitleFieldEnglish),
		verniy.MediaFieldDescription,
		verniy.MediaFieldStatus,
		verniy.MediaFieldAverageScore,
		verniy.MediaFieldSiteURL,
		verniy.MediaFieldCoverImage(verniy.MediaCoverImageFieldColor),
		verniy.MediaFieldCoverImage(verniy.MediaCoverImageFieldExtraLarge),
	)
	if err != nil {
		return nil, err
	}

	if len(res.Media) == 0 {
		return nil, nil
	}

	return &res.Media[0], nil
}

func DisplayAnime(anime *verniy.Media) discord.MessageCreate {

	markdownDesc, err := htmltomarkdown.ConvertString(*anime.Description)
	if err != nil {
		markdownDesc = *anime.Description
	}

	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewMediaGallery(discord.MediaGalleryItem{Media: discord.UnfurledMediaItem{URL: *anime.BannerImage}}),
			discord.NewSection(discord.NewTextDisplayf("## %s", *anime.Title.English)).
				WithAccessory(discord.NewSecondaryButton(string(*anime.Status), "anime:statusbutton").WithDisabled(true)),
			discord.NewSection(discord.NewTextDisplay(markdownDesc)).
				WithAccessory(discord.NewSecondaryButton("⭐ "+strconv.Itoa(*anime.AverageScore)+"/100", "anime_scorebutton").WithDisabled(true)),
		).WithAccentColor(0x0098FF),
		discord.NewActionRow(
			discord.NewLinkButton("AniList", *anime.SiteURL),
		),
	)
}
