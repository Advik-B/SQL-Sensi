package management

// Import timestamp package
import (
	"time"
)

// import telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"



type User struct {
	ID             int64
	Username       string
	FName          string
	LName          string
	LanguageCode   string
	CreatedAt      time.Time
	IsAdmin        bool
	SQLUsername    string
	SQLPassword    string
	SQLDBName      string
	GeminiAPIKey   string
}

