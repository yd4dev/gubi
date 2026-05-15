package shared

import (
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/yd4dev/gubi/internal/bot"
	"github.com/yd4dev/gubi/internal/database"
)

func ViewCheckList(owner snowflake.ID, name string, page int) discord.MessageCreate {

	var checklist database.Checklist
	if err := database.DB.Where("owner = ? AND name = ?", owner, name).Preload("Entries").First(&checklist).Error; err != nil {
		return bot.ErrorMessage(fmt.Sprintf("You don't have a checklist named '%s'.", name)).WithEphemeral(true)
	}

	limit := 11

	entries := checklist.Entries[(page-1)*limit : min(page*limit, len(checklist.Entries))]

	container := discord.NewContainer(
		discord.NewTextDisplayf("### %s (Page %d)", checklist.Name, page),
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
				discord.NewSecondaryButton("Edit", "list:editentry:"+strconv.FormatUint(uint64(entry.ID), 10)+":"+strconv.FormatInt(int64(page), 10)).WithEmoji(discord.NewComponentEmoji("✏️")),
			),
		)
	}

	container = container.AddComponents(
		discord.NewActionRow(
			discord.NewSecondaryButton("⬅️", "list:page:"+strconv.FormatUint(uint64(checklist.ID), 10)+":"+strconv.FormatInt(int64(page-1), 10)).WithDisabled(page <= 1),
			discord.NewSuccessButton("Add Entry", "list:addentry:"+strconv.FormatUint(uint64(checklist.ID), 10)+":"+strconv.FormatInt(int64(page), 10)).WithEmoji(discord.NewComponentEmoji("➕")),
			discord.NewSecondaryButton("➡️", "list:page:"+strconv.FormatUint(uint64(checklist.ID), 10)+":"+strconv.FormatInt(int64(page+1), 10)).WithDisabled(page*limit >= len(checklist.Entries)),
		),
	)

	return discord.NewMessageCreateV2(
		container,
	)
}

func ViewCheckListByID(owner snowflake.ID, checklistID uint, page int) discord.MessageCreate {

	var checklist database.Checklist
	if err := database.DB.Where("owner = ? AND id = ?", owner, checklistID).Preload("Entries").First(&checklist).Error; err != nil {
		return bot.ErrorMessage("Checklist not found.").WithEphemeral(true)
	}

	return ViewCheckList(owner, checklist.Name, page)
}
