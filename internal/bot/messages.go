package bot

import (
	"github.com/disgoorg/disgo/discord"
)

func ErrorMessage(text string) discord.MessageCreate {
	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay(text),
		).WithAccentColor(0xFF1F43),
	)
}

func SuccessMessage(text string) discord.MessageCreate {
	return discord.NewMessageCreateV2(
		discord.NewContainer(
			discord.NewTextDisplay(text),
		).WithAccentColor(0x32CD32),
	)
}
