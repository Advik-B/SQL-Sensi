package database

import (
	"fmt"
	"log"
)

func (m *MySQL) CreateUser(username string, password string) error {
	// Log the action
	log.Printf("Creating user %s", username)
	// Construct the query using the escaped values
	query := fmt.Sprintf("CREATE USER '%s'@'%%' IDENTIFIED BY '%s'", username, password)

	// Execute the query
	_, err := m.Conn.Exec(query)
	if err != nil {
		log.Print(err)
		log.Println("Possible cause: The telegram.user table is dropped but the user is not removed from MySQL users")
		return err
	}
	return nil
}