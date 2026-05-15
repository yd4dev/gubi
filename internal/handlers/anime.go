package handlers

import (
	"errors"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/yd4dev/gubi/internal/bot"
	"github.com/yd4dev/gubi/internal/shared"
)

func init() {
	bot.RegisterInteractionHandler("anime", AnimeHandler)
}

func AnimeHandler(event *events.ComponentInteractionCreate) error {
	split := strings.Split(event.Data.CustomID(), ":")

	if len(split) != 3 {
		return errors.New("invalid customID format")
	}

	switch split[1] {
	case "search":
		{
			name := split[2]
			anime, err := shared.SearchAnimeByName(name)

			if err != nil {
				return err
			}

			if anime == nil {
				return event.CreateMessage(discord.NewMessageCreate().WithContentf("No anime found with name %s", name).WithEphemeral(true))
			}

			return event.CreateMessage(shared.DisplayAnime(anime))
		}
	}
	return nil
}
