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
	// Check if the user exists in the database
	if !u.ExistsInDataBase(db) {
		u.AddToDataBase(db)
	} else {
		u.GetFromDataBase(db)
	}
	return u
}

func (u *User) GetFromDataBase(db *database.MySQL) {
	query := "SELECT * FROM users WHERE id = ?"
	rows, err := db.Conn.Query(query, u.ID)
	if err != nil {
		panic(err)
	}
	if !rows.Next() {
		panic("User not found in database")
	}
	err = rows.Scan(&u.ID, &u.Username, &u.FName, &u.LName, &u.LanguageCode, &u.SQLUsername, &u.SQLPassword, &u.SQLDBName)
	if err != nil {
		panic(err)
	}
}

func (u *User) AddToDataBase(db *database.MySQL) {
	// Convert the ID->string->byte->hash
	password, err := bcrypt.GenerateFromPassword([]byte(strconv.FormatInt(u.ID, 10)), bcrypt.DefaultCost)
	if err != nil {
	    panic(err)
	}
	u.SQLPassword = string(password)
	u.SQLDBName = "user_" + strconv.FormatInt(u.ID, 10)
	u.Username = "u" + strconv.FormatInt(u.ID, 10)
	query := "INSERT INTO users (id, username, first_name, last_name, language_code, sql_username, sql_password, sql_db_name) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err = db.Conn.Exec(query, u.ID, u.Username, u.FName, u.LName, u.LanguageCode, u.SQLUsername, u.SQLPassword, u.SQLDBName)
	if err != nil {
		panic(err)
	}
}

func (u *User) ExistsInDataBase(db *database.MySQL) bool {
	query := "SELECT * FROM users WHERE id = ?"
	rows, err := db.Conn.Query(query, u.ID)
	if err != nil {
		panic(err)
	}
	if !rows.Next() {
		return false
	}
	return true
}