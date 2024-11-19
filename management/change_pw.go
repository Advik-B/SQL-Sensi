package management

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Advik-B/SQL-Sensi/database"
	"golang.org/x/crypto/bcrypt"
)

func ChangePassword(user *User, newPassword string, db *database.MySQL) error {
	if db.Conn == nil {
		return fmt.Errorf("database connection is nil")
	}
	db.UseDatabase("telegram")
	query := `
		UPDATE users 
		SET sql_password = ?
		WHERE id = ?
	`
	_, err := db.Conn.Exec(query, newPassword, user.ID)
	if err != nil {
		return err
	}

	query = fmt.Sprintf("ALTER USER '%s'@'%%' IDENTIFIED BY '%s'", user.SQLUsername, newPassword)
	_, err = db.Conn.Exec(query)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec("FLUSH PRIVILEGES")
	return nil
}

func ResetPassword(user *User, db *database.MySQL) error {
	if db.Conn == nil {
		return fmt.Errorf("database connection is nil")
	}
	newpw, err := bcrypt.GenerateFromPassword([]byte(strconv.Itoa(time.Now().Nanosecond())), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	db.UseDatabase("telegram")
	query := `
		UPDATE users 
		SET sql_password = ?
		WHERE id = ?
	`
	_, err = db.Conn.Exec(query, string(newpw)[:20], user.ID)
	if err != nil {
		return err
	}

	query = fmt.Sprintf("ALTER USER '%s'@'%%' IDENTIFIED BY ''", user.SQLUsername)
	_, err = db.Conn.Exec(query)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec("FLUSH PRIVILEGES")
	return nil
}
