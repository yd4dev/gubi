package bot

import (
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var registeredHandlers map[string]Handler = make(map[string]Handler)

type Handler func(event *events.ComponentInteractionCreate) error

func RegisterInteractionHandler(prefix string, handler Handler) {
	registeredHandlers[prefix] = handler
}

func ComponentInteractionHandler(event *events.ComponentInteractionCreate) {
	for prefix, handler := range registeredHandlers {
		if strings.HasPrefix(event.Data.CustomID(), prefix) {
			if err := handler(event); err != nil {
				slog.Error("Error occured while handling component interaction.", slog.String("eventID", event.Data.CustomID()), slog.Any("err", err))
				event.CreateMessage(discord.NewMessageCreate().WithContent("An error occured.").WithEphemeral(true))
			}
			return
		}
	}
	slog.Error("No handler was found for component interaction.", slog.String("eventID", event.Data.CustomID()))
	event.CreateMessage(discord.NewMessageCreate().WithContent("An error occured.").WithEphemeral(true))
}
