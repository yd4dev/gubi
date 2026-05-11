package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	gubiBot "github.com/yd4dev/gubi/internal/bot"
	_ "github.com/yd4dev/gubi/internal/commands"
	"github.com/yd4dev/gubi/internal/config"
	"github.com/yd4dev/gubi/internal/database"
	_ "github.com/yd4dev/gubi/internal/handlers"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

func main() {
	slog.Info("Starting Bot...")
	slog.Info("Disgo Version", slog.String("version", disgo.Version))

	slog.Info("Loading configuration...")
	config := config.Load()

	slog.Info("Initializing database...")
	if err := database.InitDB(); err != nil {
		slog.Error("There was an error initializing the database. Exiting.", slog.Any("err", err))
		return
	}

	token := config.Token
	if token == "" {
		slog.Error("BOT_TOKEN environment variable is not set. Exiting.")
		return
	}

	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
		),
		bot.WithEventListeners(&events.ListenerAdapter{
			OnApplicationCommandInteraction: gubiBot.CommandHandler,
			OnComponentInteraction:          gubiBot.ComponentInteractionHandler,
		}),
	)
	if err != nil {
		slog.Error("Error while building disgo. Exiting.", slog.Any("err", err))
		return
	}

	defer client.Close(context.Background())

	if err = client.OpenGateway(context.Background()); err != nil {
		slog.Error("Error while connecting to gateway. Exiting.", slog.Any("err", err))
		return
	}

	if config.DevGuild != 0 {
		slog.Info("Registering Guild Commands...", slog.Uint64("guildID", uint64(config.DevGuild)))

		if _, err = client.Rest.SetGuildCommands(client.ApplicationID, config.DevGuild, gubiBot.DiscordApplicationCommandCreates); err != nil {
			slog.Error("Error while registering Guild Commands. Exiting.", slog.Any("err", err))
			return
		}
	} else {
		slog.Info("Registering Global Commands...")

		if _, err = client.Rest.SetGlobalCommands(client.ApplicationID, gubiBot.DiscordApplicationCommandCreates); err != nil {
			slog.Error("Error while registering Global Commands. Exiting.", slog.Any("err", err))
			return
		}
	}
	slog.Info("Successfully registered commands.", slog.Int("amount", len(gubiBot.RegisteredCommands)))

	selfUser, _ := client.Caches.SelfUser()

	slog.Info("Bot is now running. Press CTRL-C to exit.", slog.String("username", selfUser.Username))
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
