package main

import (
	"log"
	"os"
	"sql.sensi/commands"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	env "github.com/joho/godotenv"
)

func main() {
	env.Load() // Load .env file if it exists
	token := os.Getenv("TELEGRAM_TOKEN")
	bot, err := telegram.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	commands.RegisterCommand(commands.Command{
		Name: "start",
		Description: "Start the bot for the first time",
		Handler: func(bot *telegram.BotAPI, message *telegram.Message) {
			msg := telegram.NewMessage(message.Chat.ID, "Hello! I'm a bot that can help you with your daily tasks. Use /help to see all available commands.")
			bot.Send(msg)
		},
		Usage: "/start",
	})


	for update := range updates {
		switch {
		case update.Message == nil:
			continue
		case update.Message.IsCommand():
			for _, command := range commands.Commands {
				if update.Message.Command() == command.Name {
					command.Invoke(bot, update.Message)
				}
			}
		}
	}
}