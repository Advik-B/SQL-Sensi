package database

import (
	"log"
)

func (m *MySQL) CreateUser(username string, password string) error {

	/*
	python:
	cursor.execute(
                f"CREATE USER IF NOT EXISTS '{sql_username}'@'%' IDENTIFIED BY '{password}'"
            )
	*/

	// Create a user
	_, err := m.Conn.Exec("CREATE USER " + username + "@'%' IDENTIFIED BY '" + password + "'")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}