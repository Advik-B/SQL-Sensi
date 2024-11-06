package management

import (
	"sql.sensi/database"
)

func PrepareDB(db *database.MySQL) {
	// Connect to the MySQL database
	err := db.Connect()
	if err != nil {
		panic(err)
	}

	db.CreateAndUseDB("telegram")

	// Create the table
	err = db.CreateTable("users", []string{
		"id INT AUTO_INCREMENT PRIMARY KEY",
		"username VARCHAR(32)",
		"first_name VARCHAR(32)",
		"last_name VARCHAR(32)",
		"language_code VARCHAR(8)",
		"created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP",
		"is_admin BOOLEAN DEFAULT FALSE",
		"is_premium BOOLEAN DEFAULT FALSE",
		"sql_username VARCHAR(32) NOT NULL",
		"sql_password VARCHAR(32) NOT NULL",
		"sql_db_name VARCHAR(32) NOT NULL",
		"gemini_api_key VARCHAR(64)",
	})
	if err != nil {
		panic(err)
	}
}

