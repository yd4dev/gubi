package commands

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/yd4dev/gubi/internal/bot"
	"github.com/yd4dev/gubi/internal/database"
)

func init() {
	bot.RegisterCommand(&ListCommand{})
}

type ListCommand struct{}

func (*ListCommand) Definition() discord.ApplicationCommandCreate {
	return discord.SlashCommandCreate{
		Name:        "list",
		Description: "Manage your checklists!",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionSubCommand{
				Name:        "create",
				Description: "Create a new checklist.",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:        "name",
						Description: "The name of the checklist you want to create.",
						Required:    true,
					},
				},
			},
			discord.ApplicationCommandOptionSubCommand{
				Name:        "view",
				Description: "View a checklist.",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:         "name",
						Description:  "The name of the checklist you want to view.",
						Required:     true,
						Autocomplete: true,
					},
				},
			},
		},
		IntegrationTypes: []discord.ApplicationIntegrationType{
			discord.ApplicationIntegrationTypeGuildInstall,
			discord.ApplicationIntegrationTypeUserInstall,
		},
		Contexts: []discord.InteractionContextType{
			discord.InteractionContextTypeBotDM,
			discord.InteractionContextTypeGuild,
			discord.InteractionContextTypePrivateChannel,
		},
	}
}

func (*ListCommand) Handler(event *events.ApplicationCommandInteractionCreate) error {
	subCommand := *event.SlashCommandInteractionData().SubCommandName

	switch subCommand {
	case "create":
		return create(event)
	case "view":
		return view(event)
	}
	return errors.New("Subcommand " + subCommand + " not found.")
}

func create(event *events.ApplicationCommandInteractionCreate) error {
	name, _ := event.SlashCommandInteractionData().OptString("name")

	var checklist database.Checklist
	if err := database.DB.Where("owner = ? AND name = ?", event.User().ID, name).First(&checklist).Error; err == nil {
		return event.CreateMessage(bot.ErrorMessage(fmt.Sprintf("You already have a checklist named '%s'.", name)).WithEphemeral(true))
	}

	checklist = database.Checklist{
		Owner: event.User().ID,
		Name:  name,
	}

	if err := database.DB.Create(&checklist).Error; err != nil {
		return err
	}

	return event.CreateMessage(bot.SuccessMessage(fmt.Sprintf("Checklist '%s' created successfully!", name)).WithEphemeral(true))
}

func view(event *events.ApplicationCommandInteractionCreate) error {
	name, _ := event.SlashCommandInteractionData().OptString("name")

	var checklist database.Checklist
	if err := database.DB.Where("owner = ? AND name = ?", event.User().ID, name).Preload("Entries").First(&checklist).Error; err != nil {
		return event.CreateMessage(bot.ErrorMessage(fmt.Sprintf("You don't have a checklist named '%s'.", name)).WithEphemeral(true))
	}

	limit := 10

	entries := checklist.Entries[:min(limit, len(checklist.Entries))]

	container := discord.NewContainer(
		discord.NewTextDisplayf("### %s", checklist.Name),
	)

	for _, entry := range entries {
		var status string
		if entry.IsCompleted {
			status = "☑"
		} else {
			status = "☐"
		}
		container = container.AddComponents(
			discord.NewSection(
				discord.NewTextDisplayf("%s %s", status, entry.Description),
			).WithAccessory(
				discord.NewSecondaryButton("Edit", "list_editentry_"+strconv.FormatUint(uint64(entry.ID), 10)).WithEmoji(discord.NewComponentEmoji("✏️")),
			),
		)
	}

	container = container.AddComponents(
		discord.NewActionRow(
			discord.NewSuccessButton("Add Entry", "list_addentry_"+strconv.FormatUint(uint64(checklist.ID), 10)).WithEmoji(discord.NewComponentEmoji("➕")),
		),
	)

	return event.CreateMessage(
		discord.NewMessageCreateV2(
			container,
		),
	)
}
