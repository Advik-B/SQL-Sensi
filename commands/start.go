package commands

import (
	"fmt"
	"log"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sql.sensi/management"
)

var WelcomeBackMessage = `
Hello %s!, welcome back to SQL Sensi, your account already exists!
`
var WelcomeMessage = `
Hello %s!, welcome to SQL Sensi
I will help you learn and develop with MySQL

Here are some commands to get you started:
- /sample - Get a few sample tables to play with
- /sql - Run SQL queries
- /ai - Get help from the AI to solve your MySQL problems
- /start - Create your account/check the status of your account
- /welcome - Show this message again
... and many more

You can use /help to see all the commands available
`

func start(bot *telegram.BotAPI, message *telegram.Message) {
	if management.UserExists(&DB, message.From.ID) {
		bot.Send(telegram.NewMessage(message.Chat.ID, fmt.Sprintf(WelcomeBackMessage, message.From.FirstName)))
		return
	}
	management.UserFromTelegram(message.From, &DB)
	bot.Send(telegram.NewMessage(message.Chat.ID, fmt.Sprintf(WelcomeMessage, message.From.FirstName)))
	log.Printf("User %s has been added to the database", message.From.FirstName)
}

func welcome(bot *telegram.BotAPI, message *telegram.Message) {
	bot.Send(telegram.NewMessage(message.Chat.ID, fmt.Sprintf(WelcomeMessage, message.From.FirstName)))
}


func init() {
	Register(Command{
		Name:        "start",
		Description: "Start the bot and create your account",
		Handler:     start,
		Usage:       "/start",
	})
	Register(Command{
		Name:        "welcome",
		Description: "Show the welcome message again",
		Handler:     welcome,
		Usage:       "/welcome",
	})
}