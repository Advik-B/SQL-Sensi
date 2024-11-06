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

	// Test connection to the MySQL database
	err = mysql.Ping()
	if err != nil {
		panic(err)
	}

	mysql.CreateDatabase("test")
	// mysql.UseDatabase("test")
	whereami, err := mysql.Conn.Query("SELECT DATABASE()")	
	for whereami.Next() {
		var name string = "<NO DATABASE>"
		whereami.Scan(&name)
		println(name)
	}

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

	// Disconnect from the MySQL database
	mysql.Disconnect()
}