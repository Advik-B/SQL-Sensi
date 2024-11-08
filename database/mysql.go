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
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/", m.User, m.Password, m.Host, m.Port)
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

/*
	Returns the current database name

	if there is an error, it will return "<ERROR>"

	if we are not inside a database, it will return "<NO DATABASE>"
*/
func (m *MySQL) WhereAmI() string {
	whereami, err := m.Conn.Query("SELECT DATABASE()")
	if err != nil {
		log.Println(err)
		return "<ERROR>"
	}
	var name string = "<NO DATABASE>"
	for whereami.Next() {
		whereami.Scan(&name)
	}
	return name
}

func (m *MySQL) UseDatabase(name string) error {
	if m.Conn == nil {
		return fmt.Errorf("Database connection is nil")
	}
	// Use the database - Note: Database names cannot be parameterized in MySQL
	// However, we'll validate the name to prevent SQL injection
	if !isValidIdentifier(name) {
		return fmt.Errorf("invalid database name")
	}
	_, err := m.Conn.Exec("USE " + name)
	if err != nil {
		log.Println(err)
		return err
	}
	m.Database = name
	return nil
}

// CreateDatabase creates the database
func (m *MySQL) CreateDatabase(name string) error {
	if m.Conn == nil {
		return fmt.Errorf("Database connection is nil")
	}
	// Create the database - Note: Database names cannot be parameterized in MySQL
	// However, we'll validate the name to prevent SQL injection
	if !isValidIdentifier(name) {
		return fmt.Errorf("invalid database name")
	}
	_, err := m.Conn.Exec("CREATE DATABASE IF NOT EXISTS " + name)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Helper function to create and use a database in one go
func (m *MySQL) CreateAndUseDB(name string) error {
	// Create and use the database
	err := m.CreateDatabase(name)
	if err != nil {
		return err
	}
	err = m.UseDatabase(name)
	if err != nil {
		return err
	}
	return nil
}

// isValidIdentifier checks if a database/table name is valid
func isValidIdentifier(name string) bool {
	// Basic validation for MySQL identifiers
	// Only allow alphanumeric characters, underscores, and hyphens
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' ||
			char == '-') {
			return false
		}
	}
	return len(name) > 0
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
	if m.Conn == nil {
		return fmt.Errorf("Database connection is nil")
	}
	if !isValidIdentifier(name) {
		return fmt.Errorf("invalid table name")
	}
	// Create the table - Note: Table names and column definitions cannot be parameterized
	// However, we validate the identifiers
	_, err := m.Conn.Exec("CREATE TABLE IF NOT EXISTS " + name + " (" + strings.Join(columns, ", ") + ")")
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *MySQL) Insert(table string, columns []string, values []string) error {
	if m.Conn == nil {
		return fmt.Errorf("Database connection is nil")
	}
	if !isValidIdentifier(table) {
		return fmt.Errorf("invalid table name")
	}

	// Create placeholders for values (?)
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}

	query := "INSERT IGNORE INTO " + table + " (" + strings.Join(columns, ", ") + ") VALUES (" + strings.Join(placeholders, ", ") + ")"

	// Convert values to []interface{}
	args := make([]interface{}, len(values))
	for i, v := range values {
		args[i] = v
	}

	_, err := m.Conn.Exec(query, args...)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *MySQL) Select(table string, columns []string, where string) (*sql.Rows, error) {
	if m.Conn == nil {
		return nil, fmt.Errorf("Database connection is nil")
	}
	if !isValidIdentifier(table) {
		return nil, fmt.Errorf("invalid table name")
	}

	query := "SELECT " + strings.Join(columns, ", ") + " FROM " + table
	if where != "" {
		query += " WHERE " + where
	}

	rows, err := m.Conn.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return rows, nil
}

func (m *MySQL) Update(table string, columns []string, values []string, where string) error {
	if m.Conn == nil {
		return fmt.Errorf("Database connection is nil")
	}
	if !isValidIdentifier(table) {
		return fmt.Errorf("invalid table name")
	}

	// Create SET clause with placeholders
	setClause := make([]string, len(columns))
	for i, col := range columns {
		setClause[i] = col + " = ?"
	}

	query := "UPDATE " + table + " SET " + strings.Join(setClause, ", ")
	if where != "" {
		query += " WHERE " + where
	}

	// Convert values to []interface{}
	args := make([]interface{}, len(values))
	for i, v := range values {
		args[i] = v
	}

	_, err := m.Conn.Exec(query, args...)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *MySQL) Delete(table string, where string) error {
	if m.Conn == nil {
		return fmt.Errorf("Database connection is nil")
	}
	if !isValidIdentifier(table) {
		return fmt.Errorf("invalid table name")
	}

	query := "DELETE FROM " + table
	if where != "" {
		query += " WHERE " + where
	}

	_, err := m.Conn.Exec(query)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *MySQL) DropTable(name string) error {
	if m.Conn == nil {
		return fmt.Errorf("Database connection is nil")
	}
	if !isValidIdentifier(name) {
		return fmt.Errorf("invalid table name")
	}

	_, err := m.Conn.Exec("DROP TABLE IF EXISTS " + name)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *MySQL) DropDatabase(name string) error {
	if m.Conn == nil {
		return fmt.Errorf("Database connection is nil")
	}
	if !isValidIdentifier(name) {
		return fmt.Errorf("invalid database name")
	}

	_, err := m.Conn.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *MySQL) TruncateTable(name string) error {
	if m.Conn == nil {
		return fmt.Errorf("Database connection is nil")
	}
	if !isValidIdentifier(name) {
		return fmt.Errorf("invalid table name")
	}

	_, err := m.Conn.Exec("TRUNCATE TABLE " + name)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *MySQL) ShowDatabases() (*sql.Rows, error) {
	if m.Conn == nil {
		return nil, fmt.Errorf("Database connection is nil")
	}
	rows, err := m.Conn.Query("SHOW DATABASES")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return rows, nil
}

func (m *MySQL) ShowTables() (*sql.Rows, error) {
	if m.Conn == nil {
		return nil, fmt.Errorf("Database connection is nil")
	}
	rows, err := m.Conn.Query("SHOW TABLES")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return rows, nil
}

func (m *MySQL) DescribeTable(name string) (*sql.Rows, error) {
	if m.Conn == nil {
		return nil, fmt.Errorf("Database connection is nil")
	}
	if !isValidIdentifier(name) {
		return nil, fmt.Errorf("invalid table name")
	}

	rows, err := m.Conn.Query("DESCRIBE " + name)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return rows, nil
}

func (m *MySQL) Query(query string) (*sql.Rows, error) {
	if m.Conn == nil {
		return nil, fmt.Errorf("Database connection is nil")
	}
	rows, err := m.Conn.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return rows, nil
}

func (m *MySQL) Ping() error {
	err := m.Conn.Ping()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *MySQL) Connection() *sql.DB {
	return m.Conn
}
