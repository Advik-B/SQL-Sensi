package main

import (
	"fmt"
	"github.com/Advik-B/SQL-Sensi/commands"
	"github.com/Advik-B/SQL-Sensi/database"
	"github.com/Advik-B/SQL-Sensi/management"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	env "github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	env.Load() // Load .env file if it exists
	token := os.Getenv("TELEGRAM_TOKEN")
	bot, err := telegram.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	// bot.Debug = true
	db, err := database.FromEnvironment()
	if err != nil {
		log.Panic(err)
	}
	defer db.Disconnect()
	management.PrepareDB(db)
	commands.Initialize(db)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Start HTTP server in a separate goroutine
	go startHTTPServer()

	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		switch {
		case update.CallbackQuery != nil:
			go commands.HandleCallback(bot, &update) // Pass the callback query to the callback handler
		case update.Message != nil && update.Message.IsCommand():
			go commands.Handle(bot, update.Message) // Pass the message to the command handler
		default:
			log.Printf("Received a non-command message %v", update)
		}
	}
}

func startHTTPServer() {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Bot is running and operational.")
	})

	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Panicf("Failed to start HTTP server: %v", err)
	}
}
