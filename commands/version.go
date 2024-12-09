package commands

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func init() {
	Register(Command{
		Name:        "version",
		Description: "Get the bot version",
		Handler: func(bot *telegram.BotAPI, message *telegram.Message) {
			msg := telegram.NewMessage(message.Chat.ID, "I dont know what version I am!")
			bot.Send(msg)
		},
		Usage: "/version",
	})
}
