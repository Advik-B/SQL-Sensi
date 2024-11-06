package main

import (
	"log"
	"os"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	env "github.com/joho/godotenv"
	"sql.sensi/commands"
	"sql.sensi/database"
	"sql.sensi/management"
)

func main() {
	env.Load() // Load .env file if it exists
	token := os.Getenv("TELEGRAM_TOKEN")
	bot, err := telegram.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	db, err := database.FromEnvironment()
	if err != nil {
		log.Panic(err)
	}
	management.PrepareDB(db)
	commands.Initialize(db)

	log.Printf("Authorized on account %s", bot.Self.UserName)


	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		switch {
		case update.Message == nil:
			continue
		case update.Message.IsCommand():
			go commands.Handle(bot, update.Message) // Pass the message to the command handler
		}
	}
}
