package database

// import (
// 	"database/sql"
// 	"fmt"
// )

func New(host, user, password string) (*MySQL, error) {
	// Return the MySQL instance
	return &MySQL{
		Host:     host,
		Port:     "3306",
		User:     user,
		Password: password,
	}, nil
}
