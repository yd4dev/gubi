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
		verniy.MediaFieldTitle(verniy.MediaTitleFieldRomaji),
		verniy.MediaFieldTitle(verniy.MediaTitleFieldNative),
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

	var description = "No Description"
	if anime.Description != nil {
		description = *anime.Description
	}

	markdownDesc, err := htmltomarkdown.ConvertString(description)
	if err != nil {
		markdownDesc = description
	}

	var title = "No Title"

	if anime.Title.English != nil {
		title = *anime.Title.English
	} else if anime.Title.Romaji != nil {
		title = *anime.Title.Romaji
	} else if anime.Title.Native != nil {
		title = *anime.Title.Native
	}

	container := discord.NewContainer()
	if anime.BannerImage != nil {
		container = container.AddComponents(
			discord.NewMediaGallery(discord.MediaGalleryItem{Media: discord.UnfurledMediaItem{URL: *anime.BannerImage}}),
		)
	}

	var status = "Unknown"
	if anime.Status != nil {
		status = string(*anime.Status)
	}

	var score = "?"
	if anime.AverageScore != nil {
		score = strconv.Itoa(*anime.AverageScore)
	}

	container = container.AddComponents(
		discord.NewSection(discord.NewTextDisplayf("## %s", title)).
			WithAccessory(discord.NewSecondaryButton(status, "anime:statusbutton").WithDisabled(true)),
		discord.NewSection(discord.NewTextDisplay(markdownDesc)).
			WithAccessory(discord.NewSecondaryButton("⭐ "+score+"/100", "anime_scorebutton").WithDisabled(true)),
	).WithAccentColor(0x0098FF)

	message := discord.NewMessageCreateV2(
		container,
	)

	if anime.SiteURL != nil {
		message = message.AddComponents(
			discord.NewActionRow(
				discord.NewLinkButton("AniList", *anime.SiteURL),
			),
		)
	}
	return message
}
