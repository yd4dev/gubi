package handlers

import (
	"errors"
	"strings"

	"github.com/disgoorg/disgo/events"
	"github.com/yd4dev/gubi/internal/bot"
	"github.com/yd4dev/gubi/internal/shared"
)

func init() {
	bot.RegisterInteractionHandler("twitter", TwitterHandler)
}

func TwitterHandler(event *events.ComponentInteractionCreate) error {
	split := strings.Split(event.Data.CustomID(), ":")

	if len(split) != 2 {
		return errors.New("invalid customID format")
	}

	tweetID := split[1]

	event.DeferCreateMessage(true)

	return shared.FetchAndDisplayTweet(event.Client(), event.Token(), tweetID)
}
