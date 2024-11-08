package commands

import (
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// showTyping shows the typing indicator and returns a function to stop it
func showTyping(bot *telegram.BotAPI, chatID int64) func() {
	timer := time.AfterFunc(time.Second, func() {
		bot.Send(telegram.NewChatAction(chatID, telegram.ChatTyping))
	})
	return timer.Stop
}
