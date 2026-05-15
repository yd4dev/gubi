package bot

import (
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var registeredModalHandlers map[string]ModalHandler = make(map[string]ModalHandler)

type ModalHandler func(event *events.ModalSubmitInteractionCreate) error

func RegisterModalInteractionHandler(prefix string, handler ModalHandler) {
	registeredModalHandlers[prefix] = handler
}

func ModalInteractionHandler(event *events.ModalSubmitInteractionCreate) {
	for prefix, handler := range registeredModalHandlers {
		if strings.HasPrefix(event.Data.CustomID, prefix) {
			if err := handler(event); err != nil {
				slog.Error("Error occured while handling modal interaction.", slog.String("eventID", event.Data.CustomID), slog.Any("err", err))
				event.CreateMessage(discord.NewMessageCreate().WithContent("An error occured.").WithEphemeral(true))
			}
			return
		}
	}
	slog.Error("No handler was found for modal interaction.", slog.String("eventID", event.Data.CustomID))
	event.CreateMessage(discord.NewMessageCreate().WithContent("An error occured.").WithEphemeral(true))
}
