package management

import (
	"github.com/Advik-B/SQL-Sensi/database"
)

func PrepareDB(db *database.MySQL) {
	// Connect to the MySQL database
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Create the database
	db.CreateAndUseDB("telegram")

	// Create the table
	err = db.CreateTable("users", []string{
		"id BIGINT UNSIGNED PRIMARY KEY",
		"username VARCHAR(32)",
		"first_name VARCHAR(32)",
		"last_name VARCHAR(32)",
		"language_code VARCHAR(8)",
		"created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP",
		"is_admin BOOLEAN DEFAULT FALSE",
		"sql_username VARCHAR(20) NOT NULL",
		"sql_password VARCHAR(20) NOT NULL",
		"sql_db_name VARCHAR(20) NOT NULL",
		"gemini_api_key VARCHAR(64)",
	})
	if err != nil {
		panic(err)
	}
	db.Disconnect()

	db.Database = "telegram" // Set the database to the telegram database
	db.Connect()             // Reconnect to the telegram database
}
