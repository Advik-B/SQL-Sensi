package commands

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func RollCommand() Command {
	return Command{
		Name:        "roll",
		Description: "Roll a dice",
		Handler: func(bot *telegram.BotAPI, message *telegram.Message) {
			msg := telegram.NewMessage(message.Chat.ID, "Rolling a dice...")
			bot.Send(msg)
			bot.Send(telegram.NewDice(message.Chat.ID))
		},
		Usage: "/roll",
	}
}

func init() {
	Register(RollCommand())
}
