package commands

import (
	"fmt"
	"log"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


func help(bot *telegram.BotAPI, message *telegram.Message) {
	msg := telegram.NewMessage(message.Chat.ID, "")
	if len(message.CommandArguments()) == 0 {
		msg.Text = "*Available commands:*\n"
		for _, command := range Commands {
			msg.Text += "/"+command.Name + " \\- " + command.Description + "\n"
		}
	} else {
		for _, command := range Commands {
			if message.CommandArguments() == command.Name {
				msg.Text = "*"+command.Name+"*" + "\n" + command.Description + "\n" + "Usage: " + command.Usage + "\n"
				break
			}
		}
		if msg.Text == "" {
			msg.Text = fmt.Sprintf("Command %s not found", message.CommandArguments())
		}
	}
	msg.ParseMode = "MarkdownV2"
	log.Println("\n" + msg.Text)
	bot.Send(msg)
}

func HelpCommand() Command {
	return Command{
		Name: "help",
		Description: "Show all available commands",
		Handler: help,
		Usage: "/help or /help <command>",
	}
}

func init() {
	Register(HelpCommand())
}