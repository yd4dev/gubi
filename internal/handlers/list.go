package handlers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/yd4dev/gubi/internal/bot"
	"github.com/yd4dev/gubi/internal/database"
)

func init() {
	bot.RegisterInteractionHandler("list", ListHandler)
	bot.RegisterModalInteractionHandler("list", ListModalHandler)
}

func ListHandler(event *events.ComponentInteractionCreate) error {
	split := strings.Split(event.Data.CustomID(), "_")

	if len(split) != 3 {
		return errors.New("invalid customID format")
	}

	switch split[1] {
	case "addentry":
		{
			checklistID, err := strconv.Atoi(split[2])
			if err != nil {
				return errors.New("invalid checklist ID in customID")
			}

			checklist := database.Checklist{}

			if err := database.DB.Where("ID = ?", checklistID).First(&checklist).Error; err != nil {
				return err
			}

			if checklist.Owner != event.User().ID {
				return event.CreateMessage(bot.ErrorMessage(fmt.Sprintf("Only %s can use this button!", discord.UserMention(checklist.Owner))).WithEphemeral(true))
			}

			return event.Modal(
				discord.ModalCreate{
					CustomID: "list_addentrymodal_" + strconv.Itoa(checklistID),
					Title:    "Add Entry to " + checklist.Name,
					Components: []discord.LayoutComponent{
						discord.NewLabel("Entry Description",
							discord.NewShortTextInput("description").WithPlaceholder("Enter the description of the checklist entry").WithRequired(true),
						),
					},
				},
			)
		}
	case "editentry":
		{
			entryID, err := strconv.Atoi(split[2])
			if err != nil {
				return errors.New("invalid entry ID in customID")
			}

			entry := database.ChecklistEntry{}

			if err := database.DB.Where("ID = ?", entryID).First(&entry).Error; err != nil {
				return err
			}

			checklist := database.Checklist{}

			if err := database.DB.Where("ID = ?", entry.ChecklistID).First(&checklist).Error; err != nil {
				return err
			}

			if checklist.Owner != event.User().ID {
				return event.CreateMessage(bot.ErrorMessage(fmt.Sprintf("Only %s can use this button!", discord.UserMention(checklist.Owner))).WithEphemeral(true))
			}

			return event.Modal(
				discord.ModalCreate{
					CustomID: "list_editentrymodal_" + strconv.Itoa(entryID),
					Title:    "Edit Entry in " + checklist.Name,
					Components: []discord.LayoutComponent{
						discord.NewLabel("Entry Description",
							discord.NewShortTextInput("description").WithPlaceholder("Enter the description of the checklist entry").WithRequired(true).WithValue(entry.Description),
						),
						discord.NewLabel("Completed", discord.CheckboxComponent{
							CustomID: "completed",
							Default:  entry.IsCompleted,
						}),
					},
				},
			)
		}
	}
	return errors.New("Subcommand " + split[1] + " not found.")
}

func ListModalHandler(event *events.ModalSubmitInteractionCreate) error {
	split := strings.Split(event.Data.CustomID, "_")

	if len(split) != 3 {
		return errors.New("invalid customID format")
	}

	switch split[1] {
	case "addentrymodal":
		{
			checklistID, err := strconv.Atoi(split[2])
			if err != nil {
				return errors.New("invalid checklist ID in customID")
			}

			checklist := database.Checklist{}

			if err := database.DB.Where("ID = ?", checklistID).First(&checklist).Error; err != nil {
				return err
			}

			if checklist.Owner != event.User().ID {
				return event.CreateMessage(bot.ErrorMessage(fmt.Sprintf("Only %s can add entries!", discord.UserMention(checklist.Owner))).WithEphemeral(true))
			}

			description := event.Data.Text("description")

			entry := database.ChecklistEntry{
				ChecklistID: uint(checklistID),
				Description: description,
			}

			if err := database.DB.Create(&entry).Error; err != nil {
				return err
			}

			return event.CreateMessage(bot.SuccessMessage("Entry added successfully!").WithEphemeral(true))
		}
	case "editentrymodal":
		{
			entryID, err := strconv.Atoi(split[2])
			if err != nil {
				return errors.New("invalid entry ID in customID")
			}

			entry := database.ChecklistEntry{}

			if err := database.DB.Where("ID = ?", entryID).First(&entry).Error; err != nil {
				return err
			}

			checklist := database.Checklist{}

			if err := database.DB.Where("ID = ?", entry.ChecklistID).First(&checklist).Error; err != nil {
				return err
			}

			if checklist.Owner != event.User().ID {
				return event.CreateMessage(bot.ErrorMessage(fmt.Sprintf("Only %s can edit entries!", discord.UserMention(checklist.Owner))).WithEphemeral(true))
			}

			checkbox, _ := event.Data.Component("completed")

			if checkbox == nil {
				return errors.New("completed checkbox not found in modal data")
			}

			checkboxComponent, ok := checkbox.(discord.CheckboxComponent)
			if !ok {
				return errors.New("completed component is not a checkbox")
			}

			description := event.Data.Text("description")

			entry.Description = description
			entry.IsCompleted = checkboxComponent.Value

			if err := database.DB.Save(&entry).Error; err != nil {
				return err
			}

			return event.CreateMessage(bot.SuccessMessage("Entry updated successfully!").WithEphemeral(true))
		}
	}

	return errors.New("Subcommand " + split[1] + " not found.")
}
