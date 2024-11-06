package commands

import "sql.sensi/database"

func Initialize(db *database.MySQL) {
	DB = *db // Set the global database variable
	Register(HelpCommand())
	Register(RollCommand())
}
