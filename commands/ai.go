package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"sql.sensi/management"
)

var chatHistory = make(map[int64][]*genai.Content)

func responseToString(resp *genai.GenerateContentResponse) string {
	var str string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				str += fmt.Sprintf("%s", part)
			}
		}
	}
	return str
}

func ai(bot *telegram.BotAPI, message *telegram.Message) {
	if !accountCreateReminder(bot, message) {
		return
	}
	account := management.UserFromTelegram(message.From, &DB)
	var geminiAPIKey string
	if account.GeminiAPIKey != "" {
		geminiAPIKey = account.GeminiAPIKey
	} else {
		geminiAPIKey = os.Getenv("GEMINI_API_KEY")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiAPIKey))
	if err != nil {
		msg := telegram.NewMessage(message.Chat.ID, "")
		msg.Text = "Error: " + "```" + err.Error() + "```"
		msg.ParseMode = "MarkdownV2"
		bot.Send(msg)
		return
	}

	// Get the message text
	text := message.CommandArguments()

	// If the message is empty, return
	if text == "" {
		msg := telegram.NewMessage(message.Chat.ID, "")
		msg.Text = "Please provide a text to generate AI"
		bot.Send(msg)
		return
	}

	// Generate AI
	model := client.GenerativeModel("gemini-1.5-pro")
	cs := model.StartChat()
	// If the user has a chat history, use it
	if len(chatHistory[message.Chat.ID]) > 0 {
		cs.History = chatHistory[message.Chat.ID]
	}
	res, err := cs.SendMessage(ctx, genai.Text(text))
	if err != nil {
		msg := telegram.NewMessage(message.Chat.ID, "")
		msg.Text = "Error: " + "```" + err.Error() + "```"
		msg.ParseMode = "MarkdownV2"
		bot.Send(msg)
		return
	}
	// Send the response
	msg := telegram.NewMessage(message.Chat.ID, responseToString(res))
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
	// Save the chat history
	chatHistory[message.Chat.ID] = append(chatHistory[message.Chat.ID], cs.History...)

	// Close the client
	client.Close()
}

func clearChatHistory(bot *telegram.BotAPI, message *telegram.Message) {
	if !accountCreateReminder(bot, message) {
		return
	}
	// Set the chat history to an empty array
	chatHistory[message.Chat.ID] = []*genai.Content{}
	msg := telegram.NewMessage(message.Chat.ID, "Chat history cleared")
	bot.Send(msg)
}


func SetAPIKey(bot *telegram.BotAPI, message *telegram.Message) {
	if !accountCreateReminder(bot, message) {
		return
	}
	account := management.UserFromTelegram(message.From, &DB)
	newAPIKey := message.CommandArguments()
	if strings.TrimSpace(newAPIKey) == "" {
		msg := telegram.NewMessage(message.Chat.ID, "Please provide a new API key")
		bot.Send(msg)
		return
	}
	account.GeminiAPIKey = newAPIKey
	management.UpdateUser(&account, &DB)
	msg := telegram.NewMessage(message.Chat.ID, "Your Gemini API key has been updated")
	bot.Send(msg)
}
func init() {
	Register(
		Command{
			Name: 	  "ai",
			Description: "Generate an AI response from your query",
			Handler: ai,
			Usage: "/ai <text>\nExample: `/ai Give me a command to create an employee table with id, name, and age columns`",
		},
	)
	Register(
		Command{
			Name: 	  "clear",
			Description: "Clear the chat history with the AI",
			Handler: clearChatHistory,
			Usage: "/clear",
		},
	)
}
