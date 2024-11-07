package commands

import (
	"fmt"
	"strconv"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/olekukonko/tablewriter"
	"sql.sensi/management"
	// "strings"
)

func sql(bot *telegram.BotAPI, message *telegram.Message) {
	// Check if the user exists
	if !accountCreateReminder(bot, message) {
		return
	}
	msg := telegram.NewMessage(message.Chat.ID, "")
	// Join the arguments to form a single string
	query := message.CommandArguments()
	if strings.TrimSpace(query) == "" {
		msg.Text = "Please provide a query"
		bot.Send(msg)
		return
	}
	user := management.UserFromTelegram(message.From, &DB)
	user_db := user.GetDB(&DB)
	user_db.Connect()
	defer user_db.Disconnect()
	user_db.UseDatabase(user.SQLDBName)

	user_db.Conn.Exec("SET SESSION sql_mode = 'ANSI_QUOTES'")
	// Execute the query
	rows, err := user_db.Conn.Query(query)
	if err != nil {
		msg.Text = "```\n" + err.Error() + "\n```"
		msg.ParseMode = "MarkdownV2"
		bot.Send(msg)
		return
	}
	defer rows.Close()

	// Get the column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	// Create a slice of interface{} and a slice of interface{} pointers
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	var dataIO strings.Builder
	// Create a table writer
	var table = tablewriter.NewWriter(&dataIO)
	table.SetHeader(columns)

	const maxRows = 40
	rowCount := 0
	totalRowCount := 0
	// Fetch rows and append them to the table
	for rows.Next() {
		rows.Scan(valuePtrs...)
		strValues := make([]string, len(values))
		for i, v := range values {
			if v != nil {
				switch v := v.(type) {
				case int64:
					strValues[i] = strconv.FormatInt(v, 10)
				case float64:
					strValues[i] = strconv.FormatFloat(v, 'f', -1, 64)
				case bool:
					strValues[i] = strconv.FormatBool(v)
				case []byte:
					strValues[i] = string(v)
				case string:
					strValues[i] = v
				default:
					strValues[i] = fmt.Sprintf("%v", v)
				}
			} else {
				strValues[i] = ""
			}
		}
		table.Append(strValues)
		rowCount++
		totalRowCount++
	
		if rowCount >= maxRows {
			table.Render()
			msg.Text = "```\n" + dataIO.String() + "\n```"
			msg.ParseMode = "MarkdownV2"
			bot.Send(msg)
	
			// Reset the table and dataIO for the next batch of rows
			dataIO.Reset()
			table = tablewriter.NewWriter(&dataIO)
			rowCount = 0
		}
	}
	if totalRowCount == 0 {
		msg.Text = "Query successful, but no rows returned ✅"
		bot.Send(msg)
		return
	}
	
	// Render and send any remaining rows
	if rowCount > 0 {
		table.Render()
		msg.Text = "```\n" + dataIO.String() + "\n```"
		msg.ParseMode = "MarkdownV2"
		bot.Send(msg)
	}

}

func init() {
	Register(Command{
		Name:        "sql",
		Description: "Execute a SQL query",
		Handler:     sql,
		Usage:       "/sql <query>",
	})
}
