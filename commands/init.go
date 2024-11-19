package commands

import "github.com/Advik-B/SQL-Sensi/database"

func Initialize(db *database.MySQL) {
	DB = *db // Set the global database variable
}
