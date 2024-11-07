package commands

import telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func constructCommandString(command string, message *telegram.Message) string {
	// Call the sql function with the SELECT query
	newQuery := command
	// Append the rest of the arguments to the query
	newQuery += " " + message.CommandArguments()
	// set the message text to the result of the sql function
	message.Text = newQuery
	// Call the sql function
	return newQuery
}


func select_(bot *telegram.BotAPI, message *telegram.Message) {
	message.Text = constructCommandString("SELECT", message)
	sql(bot, message)
}

func create(bot *telegram.BotAPI, message *telegram.Message) {
	message.Text = constructCommandString("CREATE", message)
	sql(bot, message)
}

func insert(bot *telegram.BotAPI, message *telegram.Message) {
	message.Text = constructCommandString("INSERT", message)
	sql(bot, message)
}

func update(bot *telegram.BotAPI, message *telegram.Message) {
	message.Text = constructCommandString("UPDATE", message)
	sql(bot, message)
}

func delete(bot *telegram.BotAPI, message *telegram.Message) {
	message.Text = constructCommandString("DELETE", message)
	sql(bot, message)
}

func drop(bot *telegram.BotAPI, message *telegram.Message) {
	message.Text = constructCommandString("DROP", message)
	sql(bot, message)
}

func alter(bot *telegram.BotAPI, message *telegram.Message) {
	message.Text = constructCommandString("ALTER", message)
	sql(bot, message)
}

func show(bot *telegram.BotAPI, message *telegram.Message) {
	message.Text = constructCommandString("SHOW", message)
	sql(bot, message)
}

func describe(bot *telegram.BotAPI, message *telegram.Message) {
	message.Text = constructCommandString("DESCRIBE", message)
	sql(bot, message)
}


func init() {
	Register(Command{
		Name:        "select",
		Description: "Run a SELECT query on the database",
		Handler:     select_,
		Usage:       "/select <query>",
	})
	Register(Command{
		Name:        "create",
		Description: "Run a CREATE query on the database",
		Handler:     create,
		Usage:       "/create <query>",
	})
	Register(Command{
		Name:        "insert",
		Description: "Run an INSERT query on the database",
		Handler:     insert,
		Usage:       "/insert <query>",
	})
	Register(Command{
		Name:        "update",
		Description: "Run an UPDATE query on the database",
		Handler:     update,
		Usage:       "/update <query>",
	})
	Register(Command{
		Name:        "delete",
		Description: "Run a DELETE query on the database",
		Handler:     delete,
		Usage:       "/delete <query>",
	})
	Register(Command{
		Name:        "drop",
		Description: "Run a DROP query on the database",
		Handler:     drop,
		Usage:       "/drop <query>",
	})
	Register(Command{
		Name:        "alter",
		Description: "Run an ALTER query on the database",
		Handler:     alter,
		Usage:       "/alter <query>",
	})
	Register(Command{
		Name:        "show",
		Description: "Run a SHOW query on the database",
		Handler:     show,
		Usage:       "/show <query>",
	})
	Register(Command{
		Name:        "describe",
		Description: "Run a DESCRIBE query on the database",
		Handler:     describe,
		Usage:       "/describe <query>",
	})
}
