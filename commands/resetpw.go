package commands

import (
	"github.com/Advik-B/SQL-Sensi/management"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func reset_password(bot *telegram.BotAPI, message *telegram.Message) {
	if !accountCreateReminder(bot, message) {
		return
	}
	account := management.UserFromTelegram(message.From, &DB)
	err := management.ResetPassword(&account, &DB)
	if err != nil {
		bot.Send(telegram.NewMessage(message.Chat.ID, "Error resetting password: "+err.Error()))
		return
	}
	bot.Send(telegram.NewMessage(message.Chat.ID, "Password reset successfully, use /credentials to view your new password"))
}

