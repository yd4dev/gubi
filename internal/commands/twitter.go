package commands

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/yd4dev/gubi/internal/bot"
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

type TweetResponse struct {
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

func (*TwitterCommand) Handler(event *events.ApplicationCommandInteractionCreate) error {
	url, _ := event.SlashCommandInteractionData().OptString("url")
	r, err := regexp.Compile(`/status/\d+`)
	if err != nil {
		return err
	}

	event.DeferCreateMessage(true)

	apiUrl := "https://api.vxtwitter.com/Twitter"

	if r.MatchString(url) {
		apiUrl += r.FindString(url)
	} else {
		_, err := event.Client().Rest.CreateFollowupMessage(event.Client().ApplicationID, event.Token(), bot.ErrorMessage("Invalid Twitter URL!"))
		return err
	}

	res, err := http.Get(apiUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	contentType := res.Header.Get("Content-Type")

	if contentType != "application/json" {
		_, err := event.Client().Rest.CreateFollowupMessage(event.Client().ApplicationID, event.Token(), bot.ErrorMessage("Failed to fetch tweet. It might be deleted, private, or the API is down."))
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

	container = container.AddComponents(
		discord.NewTextDisplay(discord.FormattedTimestampMention(tweetResponse.DateEpoch, discord.TimestampStyleShortDateShortTime)),
	)

	event.Client().Rest.CreateFollowupMessage(event.Client().ApplicationID, event.Token(), discord.NewMessageCreate().WithContent("Fetched Tweet."))
	_, err = event.Client().Rest.CreateFollowupMessage(event.Client().ApplicationID, event.Token(), discord.NewMessageCreateV2().AddComponents(container))
	return err
}
