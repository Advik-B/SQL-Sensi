package main

import (
	"sql.sensi/database"
	env "github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	env.Load()
	// Create a new MySQL instance from the environment variables
	mysql, err := database.FromEnvironment()
	if err != nil {
		panic(err)
	}

	// Connect to the MySQL database
	err = mysql.Connect()
	if err != nil {
		panic(err)
	}
	defer mysql.Disconnect()
	
	// Test connection to the MySQL database
	err = mysql.Ping()
	if err != nil {
		panic(err)
	}

	mysql.CreateDatabase("test")
	// mysql.UseDatabase("test")
	whereami := mysql.WhereAmI() // Returns the current database name
	println(whereami)

	rows, err := mysql.ShowDatabases()
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			panic(err)
		}
		println(name)
	}
}