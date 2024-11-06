package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Connect connects to the MySQL database
func (m *MySQL) Connect() error {
	// Connect to the database
	var dsn string
	if strings.TrimSpace(m.Database) == "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)", m.User, m.Password, m.Host, m.Port)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.User, m.Password, m.Host, m.Port, m.Database)
	}
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
		return err
	}

	m.Conn = conn
	return nil
}

// CreateDatabase creates the database
func (m *MySQL) CreateDatabase(name string) error {
	// Create the database
	_, err := m.Conn.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", name))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Disconnect disconnects from the MySQL database
func (m *MySQL) Disconnect() error {
	// Disconnect from the database
	err := m.Conn.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// CreateTable creates a table in the database
func (m *MySQL) CreateTable(name string, columns []string) error {
	// Create the table
	_, err := m.Conn.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", name, strings.Join(columns, ", ")))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
