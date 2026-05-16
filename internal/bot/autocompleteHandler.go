package bot

import (
	"log/slog"

	"github.com/disgoorg/disgo/events"
)

var registeredAutocompleteHandlers map[string]AutocompleteHandler = make(map[string]AutocompleteHandler)

type AutocompleteHandler func(event *events.AutocompleteInteractionCreate) error

func RegisterAutocompleteHandler(commandName string, handler AutocompleteHandler) {
	registeredAutocompleteHandlers[commandName] = handler
}

func AutocompleteInteractionHandler(event *events.AutocompleteInteractionCreate) {
	for commandName, handler := range registeredAutocompleteHandlers {
		if event.Data.CommandName == commandName {
			if err := handler(event); err != nil {
				slog.Error("Error occured while handling autocomplete interaction.", slog.String("commandName", event.Data.CommandName), slog.Any("err", err))
			}
			return
		}
	}
	slog.Error("No handler was found for autocomplete interaction.", slog.String("commandName", event.Data.CommandName))
}
