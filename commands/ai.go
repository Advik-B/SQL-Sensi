package commands

import (
	"context"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"os"
	"fmt"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func init() {
	RegisterCommand(
		Command{
			Command:     "ai",
			Description: "Generate AI",
			Function:    ai,
		},
}
