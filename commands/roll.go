package commands

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Just a random fun command to roll a dice (i wrote this to test the dice feature lol)
func init() {
	Register(Command{
		Name:        "roll",
		Description: "Roll a dice",
		Handler: func(bot *telegram.BotAPI, message *telegram.Message) {
			msg := telegram.NewMessage(message.Chat.ID, "Rolling a dice...")
			bot.Send(msg)
			bot.Send(telegram.NewDice(message.Chat.ID))
		},
		Usage: "/roll",
	})
}
