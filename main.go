package main

import (
	"log"
	"os"

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

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		switch {
		case update.Message == nil:
			continue
		case update.Message.IsCommand():
			msg := telegram.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start":
				msg.Text = "Hello! I'm a bot that can help you with your daily tasks. Use /help to see all available commands."
			case "help":
				msg.Text = "Available commands:\n\n" +
					"/start - Start the bot\n" +
					"/help - Show all available commands\n" +
					"/ping - Check if the bot is alive\n"
				
			case "ping":
				msg.Text = "Pong!"
			default:
				msg.Text = "I don't know that command"
			
			}
			bot.Send(msg)
		}
	}
}