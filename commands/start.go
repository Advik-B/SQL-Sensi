package commands

import (
	"fmt"
	"log"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sql.sensi/management"
)

var WelcomeBackMessage = "Hello %s!, welcome back to the SQL Sensi bot. You can now use the /help command to see the available commands."
var WelcomeMessage = "Hello %s!, welcome to the SQL Sensi bot. Your account and database have been created. You can now use the /help command to see the available commands."

func start(bot *telegram.BotAPI, message *telegram.Message) {
	if management.UserExists(&DB, message.From.ID) {
		bot.Send(telegram.NewMessage(message.Chat.ID, fmt.Sprintf(WelcomeBackMessage, message.From.FirstName)))
		return
	}
	management.UserFromTelegram(message.From, &DB)
	bot.Send(telegram.NewMessage(message.Chat.ID, fmt.Sprintf(WelcomeMessage, message.From.FirstName)))
	log.Printf("User %s has been added to the database", message.From.FirstName)
}


func StartCommand() Command {
	return Command{
		Name:        "start",
		Description: "Start the bot and create your account",
		Handler:     start,
		Usage:       "/start",
	}
}