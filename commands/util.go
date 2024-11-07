package commands

import telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
import "sql.sensi/management"

func accountCreateReminder(bot *telegram.BotAPI, message *telegram.Message) bool {
	if !management.UserExists(&DB, message.From.ID) {
		msg := telegram.NewMessage(message.Chat.ID, "You need to create an account first, use /start")
		bot.Send(msg)
		return false
	}
	return true
}