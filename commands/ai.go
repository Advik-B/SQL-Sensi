package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"sql.sensi/management"
)

var (
	chatHistory = make(map[int64][]*genai.Content)
	chatMutex   sync.RWMutex
)

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
	// Show the typing indicator
	var typing bool = true
	go func() {
		for {
			if typing {
				bot.Send(telegram.NewChatAction(message.Chat.ID, telegram.ChatTyping))
				log.Println("AI typing")
				time.Sleep(1 * time.Second)
			} else {
				return
			}
		}
	}()
	// bot.Send(telegram.NewChatAction(message.Chat.ID, telegram.ChatTyping))

	// Get the user's Gemini API key
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
		log.Println("Error creating client")
		log.Println(err)
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

	// Get chat history with read lock
	chatMutex.RLock()
	history := chatHistory[message.Chat.ID]
	chatMutex.RUnlock()

	// If the user has a chat history, use it
	if len(history) > 0 {
		cs.History = history
	}

	res, err := cs.SendMessage(ctx, genai.Text(text))
	if err != nil {
		log.Println("Error sending message to AI")
		log.Println(err)
		msg := telegram.NewMessage(message.Chat.ID, "")
		msg.Text = "Error: " + "```" + err.Error() + "```"
		msg.ParseMode = "MarkdownV2"
		bot.Send(msg)
		return
	}
	// Send the response
	msg := telegram.NewMessage(message.Chat.ID, "")
	msg.Text = parseMarkDown(responseToString(res))
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)

	// Save the chat history with write lock
	chatMutex.Lock()
	chatHistory[message.Chat.ID] = append(chatHistory[message.Chat.ID], cs.History...)
	chatMutex.Unlock()

	// Close the client
	client.Close()
	typing = false
	log.Println("AI response sent")
}

func clearChatHistory(bot *telegram.BotAPI, message *telegram.Message) {
	if !accountCreateReminder(bot, message) {
		return
	}
	// Set the chat history to an empty array with write lock
	chatMutex.Lock()
	chatHistory[message.Chat.ID] = []*genai.Content{}
	chatMutex.Unlock()

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
		if strings.TrimSpace(account.GeminiAPIKey) == "" {
			msg := telegram.NewMessage(message.Chat.ID, "You don't have a Gemini API key set, nothing to clear")
			bot.Send(msg)
			return
		}
		msg := telegram.NewMessage(message.Chat.ID, "Are you sure you want to clear your Gemini API key? This action is irreversible.")
		// Add a button to confirm the action (linked to the clearAPICallback function)
		msg.ReplyMarkup = telegram.NewInlineKeyboardMarkup(
			telegram.NewInlineKeyboardRow(
				telegram.NewInlineKeyboardButtonData("Yes", "clearAPI"),
				telegram.NewInlineKeyboardButtonData("No", "cancelAPI"),
			),
		)

		bot.Send(msg)
		return
	}
	account.GeminiAPIKey = newAPIKey
	management.UpdateUser(&account, &DB)
	msg := telegram.NewMessage(message.Chat.ID, "Your Gemini API key has been updated")
	bot.Send(msg)
}

func clearAPICallback(bot *telegram.BotAPI, query *telegram.CallbackQuery) {
	account := management.UserFromTelegram(query.From, &DB)
	account.GeminiAPIKey = ""
	management.UpdateUser(&account, &DB)
	msg := telegram.NewMessage(query.Message.Chat.ID, "Your Gemini API key has been cleared")
	bot.Send(msg)
	// Prevent the button from being clicked again

}

func cancelClearAPICallback(bot *telegram.BotAPI, query *telegram.CallbackQuery) {
	// Do nothing
}

func parseMarkDown(text string) string {
	var result strings.Builder
	inCode := false

	for i := 0; i < len(text); i++ {
		ch := text[i]
		switch ch {
		case '`':
			if inCode {
				result.WriteString("\\`")
			} else {
				result.WriteString("`")
			}
			inCode = !inCode
		case '\\', '_', '*', '[', ']', '(', ')', '~', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!':
			if inCode {
				if ch == '\\' {
					result.WriteString("\\\\")
				} else {
					result.WriteByte(ch)
				}
			} else {
				result.WriteByte('\\')
				result.WriteByte(ch)
			}
		default:
			result.WriteByte(ch)
		}
	}
	return result.String()
}

func init() {
	Register(
		Command{
			Name:        "ai",
			Description: "Generate an AI response from your query",
			Handler:     ai,
			Usage:       "/ai <text>\nExample: `/ai Give me a command to create an employee table with id, name, and age columns`",
		},
	)
	Register(
		Command{
			Name:        "clear",
			Description: "Clear the chat history with the AI",
			Handler:     clearChatHistory,
			Usage:       "/clear",
		},
	)
	Register(
		Command{
			Name:        "apikey",
			Description: "\\(advanced\\) Set/Clear your Gemini API key",
			Handler:     SetAPIKey,
			Usage:       "/apikey <API key> or /apikey to clear the API key",
		},
	)
	RegisterCallback(
		Callback{
			Name:    "clearAPI",
			Handler: clearAPICallback,
		},
	)
	RegisterCallback(
		Callback{
			Name:    "cancelAPI",
			Handler: cancelClearAPICallback,
		},
	)
}
