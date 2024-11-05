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

// Invoke method to execute the command
func (c Command) Invoke(bot *telegram.BotAPI, message *telegram.Message) {
	c.Handler(bot, message)
}

func (c Command) String() string {
	str := c.Name + " - " + c.Description + "\n"
	str += "Usage: " + c.Usage
	return str
}

func Register(command Command) {
	log.Printf("Registering command %v", command)
	Commands = append(Commands, command)
}

func Handle(bot *telegram.BotAPI, message *telegram.Message) {
	for _, command := range Commands {
		if message.IsCommand() && message.Command() == command.Name {
			command.Invoke(bot, message)
			return
		}
	}
}
