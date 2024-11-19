package management

import (
	"fmt"

	"github.com/Advik-B/SQL-Sensi/database"
)

func UpdateUser(user *User, db *database.MySQL) error {
	if db.Conn == nil {
		return fmt.Errorf("database connection is nil")
	}
	db.UseDatabase("telegram")
	query := `
		UPDATE users 
		SET username = ?, first_name = ?, last_name = ?, language_code = ?, is_admin = ?, sql_username = ?, sql_password = ?, sql_db_name = ?, gemini_api_key = ? 
		WHERE id = ?
	`
	_, err := db.Conn.Exec(query, user.Username, user.FName, user.LName, user.LanguageCode, user.IsAdmin, user.SQLUsername, user.SQLPassword, user.SQLDBName, user.GeminiAPIKey, user.ID)
	if err != nil {
		return err
	}
	return nil
}
