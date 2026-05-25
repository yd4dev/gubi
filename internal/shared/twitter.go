package shared

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"

	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/yd4dev/gubi/internal/bot"
)

type TweetResponse struct {
	QRT
	Qrt    QRT    `json:"qrt,omitempty"`
	QrtURL string `json:"qrtURL,omitempty"`
}

type QRT struct {
	Date           string          `json:"date"`
	DateEpoch      int64           `json:"date_epoch"`
	Hashtags       []string        `json:"hashtags"`
	Likes          int             `json:"likes"`
	MediaURLs      []string        `json:"mediaURLs"`
	MediaExtended  []ExtendedMedia `json:"media_extended"`
	Replies        int             `json:"replies"`
	Retweets       int             `json:"retweets"`
	Text           string          `json:"text"`
	TweetID        string          `json:"tweetID"`
	TweetURL       string          `json:"tweetURL"`
	UserName       string          `json:"user_name"`
	UserScreenName string          `json:"user_screen_name"`
}

type ExtendedMedia struct {
	AltText        string `json:"altText"`
	DurationMillis int    `json:"duration_millis,omitempty"`
	Size           Size   `json:"size"`
	ThumbnailURL   string `json:"thumbnail_url"`
	Type           string `json:"type"`
	URL            string `json:"url"`
}

type Size struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

func FetchAndDisplayTweet(client *disgobot.Client, token string, tweetID string) error {

	apiUrl := "https://api.vxtwitter.com/Twitter/status/" + tweetID

	res, err := http.Get(apiUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	contentType := res.Header.Get("Content-Type")

	if contentType != "application/json" {
		_, err := client.Rest.CreateFollowupMessage(client.ApplicationID, token, bot.ErrorMessage("Failed to fetch tweet. It might be deleted, private, or the API is down."))
		return err
	}

	var tweetResponse TweetResponse

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resBody, &tweetResponse)
	if err != nil {
		return err
	}

	container := discord.NewContainer()

	container = container.AddComponents(
		discord.NewSection(
			discord.NewTextDisplayf("**%s** (@%s)", tweetResponse.UserName, tweetResponse.UserScreenName),
		).WithAccessory(
			discord.NewLinkButton("Open URL", tweetResponse.TweetURL),
		),
	).WithAccentColor(0x1CA0F1)

	if b, err := regexp.MatchString(`^https://t\.co/\w+$`, tweetResponse.Text); err == nil && !b {
		container = container.AddComponents(
			discord.NewTextDisplay(tweetResponse.Text),
		)
	}

	var mediaGalleryItems []discord.MediaGalleryItem

	for _, media := range tweetResponse.MediaURLs {
		mediaGalleryItems = append(mediaGalleryItems, discord.MediaGalleryItem{
			Media: discord.UnfurledMediaItem{
				URL: media,
			},
		})
	}

	if len(mediaGalleryItems) > 0 {
		container = container.AddComponents(
			discord.NewMediaGallery(mediaGalleryItems...),
		)
	}

	if tweetResponse.QrtURL != "" {

		var comp discord.ContainerSubComponent

		var text = "> **%s** (@%s) - [Quote Retweet](%s)"

		var textDisplay = discord.NewTextDisplayf(text, tweetResponse.Qrt.UserName, tweetResponse.Qrt.UserScreenName, tweetResponse.Qrt.TweetURL)

		if b, err := regexp.MatchString(`^https://t\.co/\w+$`, tweetResponse.Qrt.Text); err == nil && !b {
			text += "\n> \n> %s"
			textDisplay = discord.NewTextDisplayf(text, tweetResponse.Qrt.UserName, tweetResponse.Qrt.UserScreenName, tweetResponse.Qrt.TweetURL, strings.ReplaceAll(tweetResponse.Qrt.Text, "\n", "\n> "))
		}

		if len(tweetResponse.Qrt.MediaExtended) > 0 {
			comp = discord.NewSection(
				textDisplay,
			).WithAccessory(discord.NewThumbnail(tweetResponse.Qrt.MediaExtended[0].ThumbnailURL))
		} else {
			comp = textDisplay
		}

		container = container.AddComponents(
			discord.NewSmallSeparator(),
			comp,
		)

	}

	var footer discord.ContainerSubComponent

	footerText := discord.FormattedTimestampMention(tweetResponse.DateEpoch, discord.TimestampStyleShortDateShortTime)

	footer = discord.NewTextDisplay(footerText)

	if tweetResponse.QrtURL != "" {
		footer = discord.NewSection(
			discord.NewTextDisplay(footerText),
		).WithAccessory(discord.NewPrimaryButton("Fetch Quoted Tweet", "twitter:"+tweetResponse.Qrt.TweetID))
	}

	container = container.AddComponents(
		footer,
	)

	client.Rest.CreateFollowupMessage(client.ApplicationID, token, discord.NewMessageCreate().WithContent("Fetched Tweet."))
	_, err = client.Rest.CreateFollowupMessage(client.ApplicationID, token, discord.NewMessageCreateV2().AddComponents(container))
	return err
}
