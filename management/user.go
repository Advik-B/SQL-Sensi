package management

import (
	"strconv"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/crypto/bcrypt"
	"sql.sensi/database"
)


func UserFromUpdate(update *telegram.Update, db* database.MySQL) User {
	u := User{}
	user := update.Message.From
	u.ID = int64(user.ID)
	u.Username = user.UserName
	u.FName = user.FirstName
	u.LName = user.LastName
	u.LanguageCode = user.LanguageCode

	db.UseDatabase("telegram")
	// Check if the user exists in the database
	query := "SELECT * FROM users WHERE id = ?"
	rows, err := db.Conn.Query(query, u.ID)
	if err != nil {
		panic(err)
	}
	// Check the number of rows returned
	if !rows.Next() {
		// User does not exist, insert the user
		insertQuery := "INSERT INTO users (id, username, first_name, last_name, language_code) VALUES (?, ?, ?, ?, ?)"
		_, err := db.Conn.Exec(insertQuery, u.ID, u.Username, u.FName, u.LName, u.LanguageCode)
		if err != nil {
			panic(err)
		}
	} else {
		// User exists, update the user
		updateQuery := "UPDATE users SET username = ?, first_name = ?, last_name = ?, language_code = ? WHERE id = ?"
		_, err := db.Conn.Exec(updateQuery, u.Username, u.FName, u.LName, u.LanguageCode, u.ID)
		if err != nil {
			panic(err)
		}
	}
	rows.Close()
	return u
}


func (u *User) AddToDataBase(db *database.MySQL) {
	db.UseDatabase("telegram")
	// Convert the ID->string->byte->hash
	password, err := bcrypt.GenerateFromPassword([]byte(strconv.FormatInt(u.ID, 10)), bcrypt.MaxCost)
	if err != nil {
	    panic(err)
	}

}