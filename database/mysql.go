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
	// Use the database
	_, err := m.Conn.Exec(fmt.Sprintf("USE %s", name))
	if err != nil {
		log.Println(err)
		return err
	}
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

func (m *MySQL) Insert(table string, columns []string, values []string) error {
	// Insert a row into the table
	_, err := m.Conn.Exec(fmt.Sprintf("INSERT IGNORE INTO %s (%s) VALUES (%s)", table, strings.Join(columns, ", "), strings.Join(values, ", ")))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *MySQL) Select(table string, columns []string, where string) (*sql.Rows, error) {
	// Select rows from the table
	rows, err := m.Conn.Query(fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(columns, ", "), table, where))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return rows, nil
}

func (m *MySQL) Update(table string, columns []string, values []string, where string) error {
	// Update rows in the table
	_, err := m.Conn.Exec(fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(columns, ", "), where))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *MySQL) Delete(table string, where string) error {
	// Delete rows from the table
	_, err := m.Conn.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s", table, where))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *MySQL) DropTable(name string) error {
	// Drop the table
	_, err := m.Conn.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", name))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *MySQL) DropDatabase(name string) error {
	// Drop the database
	_, err := m.Conn.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", name))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *MySQL) TruncateTable(name string) error {
	// Truncate the table
	_, err := m.Conn.Exec(fmt.Sprintf("TRUNCATE TABLE %s", name))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *MySQL) ShowDatabases() (*sql.Rows, error) {
	// Show the databases
	rows, err := m.Conn.Query("SHOW DATABASES")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return rows, nil
}

func (m *MySQL) ShowTables() (*sql.Rows, error) {
	// Show the tables
	rows, err := m.Conn.Query("SHOW TABLES")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return rows, nil
}

func (m *MySQL) DescribeTable(name string) (*sql.Rows, error) {
	// Describe the table
	rows, err := m.Conn.Query(fmt.Sprintf("DESCRIBE %s", name))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return rows, nil
}

func (m *MySQL) Query(query string) (*sql.Rows, error) {
	// Execute a custom query
	rows, err := m.Conn.Query(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return rows, nil
}

func (m *MySQL) Ping() error {
	// Ping the database
	err := m.Conn.Ping()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (m *MySQL) Connection() *sql.DB {
	// Return the connection
	return m.Conn
}
