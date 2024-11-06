package database

import (
	"os"
)

func FromEnvironment() (*MySQL, error) {
	// Return the MySQL instance
	return New(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
	)
}
