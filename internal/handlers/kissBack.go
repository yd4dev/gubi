package handlers

import (
	"errors"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/yd4dev/gubi/internal/bot"
	"github.com/yd4dev/gubi/internal/shared"
)

func init() {
	bot.RegisterInteractionHandler("kiss", KissHandler)
}

func KissHandler(event *events.ComponentInteractionCreate) error {
	kisser := event.Member()

	split := strings.Split(event.Data.CustomID(), "_")

	if len(split) != 3 {
		return errors.New("invalid customID format")
	}

	kisserID, err := snowflake.Parse(split[1])

	if err != nil {
		return errors.New("invalid user ID in customID")
	}

	kissedID, err := snowflake.Parse(split[2])

	if err != nil {
		return errors.New("invalid user ID in customID")
	}

	if kisserID != kisser.User.ID {
		event.CreateMessage(discord.NewMessageCreate().WithContentf("Only %s can use this button!", discord.UserMention(kisserID)).WithEphemeral(true))
		return nil
	}

	message, err := shared.Kiss(kisser.User.ID, kissedID)
	if err != nil {
		return err
	}

	if err := event.CreateMessage(message); err != nil {
		return err
	}
	return nil
}
