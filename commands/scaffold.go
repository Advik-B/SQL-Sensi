package commands

import (
	"log"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/Advik-B/SQL-Sensi/database"
)

type Command struct {
	Name        string
	Description string
	Handler     func(*telegram.BotAPI, *telegram.Message)
	Usage       string
}

type Callback struct {
	Name    string
	Handler func(*telegram.BotAPI, *telegram.CallbackQuery)
}

var Commands = []Command{}
var Callbacks = []Callback{}
var DB database.MySQL

func (c Command) String() string {
	str := c.Name + " - " + c.Description + "\n"
	str += "Usage: " + c.Usage
	return str
}

func Register(command Command) {
	log.Printf("Registering command: %s", command.Name)
	Commands = append(Commands, command)
}

func RegisterCallback(callback Callback) {
	log.Printf("Registering callback: %s", callback.Name)
	Callbacks = append(Callbacks, callback)
}

func Handle(bot *telegram.BotAPI, message *telegram.Message) {
	for _, command := range Commands {
		if message.Command() == command.Name {
			log.Printf("Invoking command /%s", command.Name)
			command.Handler(bot, message)
			return
		}
	}
	log.Printf("Command /%s not found", message.Command())
}

func HandleCallback(bot *telegram.BotAPI, update *telegram.Update) {
	if update.CallbackQuery == nil {
		log.Printf("A callback query is expected, but got nil?")
		return
	}
	for _, callback := range Callbacks {
		if update.CallbackQuery.Data == callback.Name {
			log.Printf("Invoking callback %s", callback.Name)
			callback.Handler(bot, update.CallbackQuery)
			return
		}
	}
	log.Printf("Callback %s not found", update.CallbackQuery.Data)
}
