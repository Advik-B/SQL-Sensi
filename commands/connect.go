package commands

import (
	"sql.sensi/management"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"fmt"
)

var connectMessage = "Here are the credentials to connect to your database\nHost: `%s`\nUsername: `%s`\nPassword: `%s`\nDatabase: `%s`"

var pythonCode = `
import mysql.connector as sql

db = sql.connect(
	host="%s",
	user="%s",
	password="%s",
	database="%s"
)

conn = db.cursor()

`

var connectMessagePython ="Some python code to quickstart:\n```python\n" + pythonCode + "```"

// I know this is pretty fucking jank but eh
func generateConnectMessage(host, username, password, database string) string {
	return fmt.Sprintf(connectMessagePython, host, username, password, database)
}


func connect(bot *telegram.BotAPI, message *telegram.Message) {
	if !accountCreateReminder(bot, message) {
		return
	}

	user := management.UserFromTelegram(message.From, &DB)
	msg := telegram.NewMessage(message.Chat.ID, "")
	msg.Text = generateConnectMessage(DB.Host, user.SQLUsername, user.SQLPassword, user.SQLDBName)
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
}

func credentials(bot *telegram.BotAPI, message *telegram.Message) {
	if !accountCreateReminder(bot, message) {
		return
	}

	user := management.UserFromTelegram(message.From, &DB)
	msg := telegram.NewMessage(message.Chat.ID, fmt.Sprintf(connectMessage, DB.Host, user.SQLUsername, user.SQLPassword, user.SQLDBName))
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
}

func init() {
	Register(Command{
		Name:        "connect",
		Description: "Get the credentials to connect to your database",
		Handler:     connect,
		Usage:       "/connect",
	})
	Register(Command{
		Name:        "credentials",
		Description: "Get the credentials to connect to your database",
		Handler:     credentials,
		Usage:       "/credentials",
	})
}
