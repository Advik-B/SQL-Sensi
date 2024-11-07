package commands

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sql(bot *telegram.BotAPI, message *telegram.Message) {
	// Join the arguments to form a single string
	query := message.CommandArguments()
	// Execute the query
	msg := telegram.NewMessage(message.Chat.ID, query)
	bot.Send(msg)
}

func init() {
	Register(Command{
		Name:        "sql",
		Description: "Execute a SQL query",
		Handler:     sql,
		Usage:       "/sql <query>",
	})
}