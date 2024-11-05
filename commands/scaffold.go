package commands

import (
	"log"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Command struct {
	Name        string
	Description string
	Handler     func(*telegram.BotAPI, *telegram.Message)
	Usage       string
}

var Commands = []Command{}


func (c Command) String() string {
	str := c.Name + " - " + c.Description + "\n"
	str += "Usage: " + c.Usage
	return str
}

func Register(command Command) {
	log.Printf("Registering command: %s", command.Name)
	Commands = append(Commands, command)
}

func Handle(bot *telegram.BotAPI, message *telegram.Message) {
	for _, command := range Commands {
		if message.Command() == command.Name {
			log.Printf("Invoking command %s", command.Name)
			command.Handler(bot, message)
			return
		}
	}
	log.Printf("Command %s not found", message.Command())
	bot.Send(telegram.NewMessage(message.Chat.ID, "Command not found. Use /help to see all available commands."))
}
