package management

import (
	"database/sql"
	"strconv"

	"time"

	"github.com/Advik-B/SQL-Sensi/database"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/crypto/bcrypt"
)

func UserFromTelegram(user *telegram.User, db *database.MySQL) User {
	u := User{}
	u.ID = int64(user.ID)
	u.Username = user.UserName
	u.FName = user.FirstName
	u.LName = user.LastName
	u.LanguageCode = user.LanguageCode
	// Check if the user exists in the database
	if !u.ExistsInDataBase(db) {
		u.AddToDataBase(db)
	} else {
		u.GetFromDataBase(db)
	}
	return u
}

func (u *User) GetFromDataBase(db *database.MySQL) {
	if db.Conn == nil {
		panic("Database connection is nil")
	}
	query := "SELECT * FROM users WHERE id = ?"
	rows, err := db.Conn.Query(query, u.ID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	if !rows.Next() {
		panic("User not found in database")
	}
	var createdAt []byte
	var geminiAPIKey sql.NullString
	err = rows.Scan(
		&u.ID,
		&u.Username,
		&u.FName,
		&u.LName,
		&u.LanguageCode,
		&createdAt,
		&u.IsAdmin,
		&u.SQLUsername,
		&u.SQLPassword,
		&u.SQLDBName,
		&geminiAPIKey,
	)
	if err != nil {
		panic(err)
	}
	u.CreatedAt, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		panic(err)
	}
	if geminiAPIKey.Valid {
		u.GeminiAPIKey = geminiAPIKey.String
	} else {
		u.GeminiAPIKey = ""
	}
}

func (u *User) AddToDataBase(db *database.MySQL) {
	if db.Conn == nil {
		panic("Database connection is nil")
	}
	// Convert the ID->string->byte->hash
	password, err := bcrypt.GenerateFromPassword([]byte(strconv.FormatInt(u.ID, 10)), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	u.SQLPassword = string(password)[:20] // Limit the password to 20 characters
	u.SQLDBName = "user_" + strconv.FormatInt(u.ID, 10)
	u.SQLUsername = "u" + strconv.FormatInt(u.ID, 10)
	query := "INSERT INTO users (id, username, first_name, last_name, language_code, sql_username, sql_password, sql_db_name) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = db.Conn.Exec(query, u.ID, u.Username, u.FName, u.LName, u.LanguageCode, u.SQLUsername, u.SQLPassword, u.SQLDBName)
	if err != nil {
		panic(err)
	}
	// Create the user's database
	db.CreateDatabase(u.SQLDBName)
	// Create the user's database user
	db.CreateUser(u.SQLUsername, u.SQLPassword)
	// Grant the user access to the user's database
	query = "GRANT ALL PRIVILEGES ON " + u.SQLDBName + ".* TO " + u.SQLUsername + "@'%'"
	_, err = db.Conn.Exec(query)
	if err != nil {
		panic(err)
	}
	// Flush the privileges
	_, err = db.Conn.Exec("FLUSH PRIVILEGES")
	if err != nil {
		panic(err)
	}
}

func UserExists(db *database.MySQL, id int64) bool {
	if db.Conn == nil {
		panic("Database connection is nil")
	}
	query := "SELECT id FROM users WHERE id = ?"
	rows, err := db.Conn.Query(query, id)
	if err != nil {
		panic(err)
	}
	if !rows.Next() {
		return false
	}
	return true
}

// Proxy function to check if the user exists in the database
func (u *User) ExistsInDataBase(db *database.MySQL) bool {
	if db.Conn == nil {
		panic("Database connection is nil")
	}
	return UserExists(db, u.ID)
}

func (u *User) GetDB(db *database.MySQL) database.MySQL {
	if db.Conn == nil {
		panic("Database connection is nil")
	}
	host := db.Host
	uDB, err := database.New(host, u.SQLUsername, u.SQLPassword)
	uDB.UseDatabase(u.SQLDBName)
	if err != nil {
		panic(err)
	}
	return *uDB
}
