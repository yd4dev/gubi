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
	"github.com/yd4dev/gubi/internal/shared"
)

func init() {
	bot.RegisterInteractionHandler("list", ListHandler)
	bot.RegisterModalInteractionHandler("list", ListModalHandler)
	bot.RegisterAutocompleteHandler("list", ListAutocompleteHandler)
}

func ListHandler(event *events.ComponentInteractionCreate) error {
	split := strings.Split(event.Data.CustomID(), ":")

	if len(split) != 3 && len(split) != 4 {
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

			page := split[3]

			return event.Modal(
				discord.ModalCreate{
					CustomID: "list:addentrymodal:" + strconv.Itoa(checklistID) + ":" + page,
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

			page := split[3]

			return event.Modal(
				discord.ModalCreate{
					CustomID: "list:editentrymodal:" + strconv.Itoa(entryID) + ":" + page,
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
	case "page":
		{
			checklistID, err := strconv.Atoi(split[2])
			if err != nil {
				return errors.New("invalid checklist ID in customID")
			}

			page, err := strconv.Atoi(split[3])
			if err != nil {
				return errors.New("invalid page number in customID")
			}

			return event.UpdateMessage(discord.NewMessageUpdateV2(shared.ViewCheckListByID(event.User().ID, uint(checklistID), page).Components...))
		}
	}

	return errors.New("Subcommand " + split[1] + " not found.")
}

func ListModalHandler(event *events.ModalSubmitInteractionCreate) error {
	split := strings.Split(event.Data.CustomID, ":")

	if len(split) != 3 && len(split) != 4 {
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

			page, err := strconv.Atoi(split[3])
			if err != nil {
				page = 1
			}

			return event.UpdateMessage(discord.NewMessageUpdateV2(shared.ViewCheckList(checklist.Owner, checklist.Name, page).Components...))
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

			page, err := strconv.Atoi(split[3])
			if err != nil {
				page = 1
			}

			return event.UpdateMessage(discord.NewMessageUpdateV2(shared.ViewCheckList(checklist.Owner, checklist.Name, page).Components...))
		}
	}

	return errors.New("Subcommand " + split[1] + " not found.")
}

func ListAutocompleteHandler(event *events.AutocompleteInteractionCreate) error {
	switch *event.Data.SubCommandName {
	case "view":
		{
			name := event.Data.Focused().String()

			checklists := []database.Checklist{}

			if name == "" {
				database.DB.Where("Owner = ?", event.User().ID).Find(&checklists)
			} else {
				database.DB.Where("Owner = ? AND Name LIKE ?", event.User().ID, name+"%").Find(&checklists)
			}

			result := []discord.AutocompleteChoice{}

			for _, list := range checklists {
				result = append(result, discord.AutocompleteChoiceString{Name: list.Name, Value: list.Name})
			}

			return event.AutocompleteResult(result)
		}
	}
	return errors.New("Subcommand " + *event.Data.SubCommandName + " not found.")
}
