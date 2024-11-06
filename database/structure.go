package database

import "database/sql"

// MySQL is an implementation of the Database interface
type MySQL struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string

	Conn *sql.DB
}
