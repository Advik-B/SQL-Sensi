package commands

import (
	"github.com/Advik-B/SQL-Sensi/management"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func accountCreateReminder(bot *telegram.BotAPI, message *telegram.Message) bool {
	if !management.UserExists(&DB, message.From.ID) {
		msg := telegram.NewMessage(message.Chat.ID, "You need to create an account first, use /start")
		bot.Send(msg)
		return false
	}
	return true
}
