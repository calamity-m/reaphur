package bot

import (
	"github.com/disgoorg/disgo/discord"
)

const (
	SlashHelpCommand   = "help"
	SlashGetCommand    = "get"
	MessageEditCommand = "edit msg"
)

var (
	BotCommands = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        SlashHelpCommand,
			Description: "get some help with how to use this thing",
		},
		discord.SlashCommandCreate{
			Name:        SlashGetCommand,
			Description: "not implemeted soz.soz",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
					Name:        "days",
					Description: "days since",
					Required:    false,
				},
			},
		},
		discord.MessageCommandCreate{
			Name: MessageEditCommand,
		},
	}
)
