package commands

import (
	"fmt"
	"sql.sensi/management"
	"log"
	"strings"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var SampleQueries = `
CREATE TABLE IF NOT EXISTS Employee (
	ID INT PRIMARY KEY AUTO_INCREMENT,
	Name VARCHAR(255) NOT NULL,
	Age INT,
	DepartNo INT,
	Salary INT NOT NULL,
	CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	Commision INT
);

CREATE TABLE IF NOT EXISTS Department (
	DepartNo INT PRIMARY KEY AUTO_INCREMENT,
	DepartName VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS Student (
	ID INT PRIMARY KEY AUTO_INCREMENT,
	Name VARCHAR(255) NOT NULL,
	Age INT,
	Grade INT,
	CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

var SampleInserts = `
INSERT IGNORE INTO Employee (Name, Age, DepartNo, Salary, Commision) VALUES 
('Roop', 32, 5, 75000, 3000),
('Adithya', 27, 5, 65000, 2000),
('Advik', 23, 5, 69696969, 1000),
('Pratibha', 29, 6, 42000, 500),
('Ashvitha', 24, 7, 42069, 1000),
('Dhruv', 30, 7, 420420, 2000),
('Dhariya Jaishwal', 35, 8, 420420420, 3000),
('Chatur', 40, 8, 420, 4000),
('Sowlaasya', 20, 9, 2040, 500),
('Geetika', 45, 9, 30000, 6000),
('Pooja', 22, 10, 40000, 1000),
('Prachi', 28, 10, 60000, 2000),
('Sakshi', 35, 11, 70000, 3000),
('John Doe', 30, 1, 50000, 1000),
('Jane Doe', 25, 1, 45000, 500),
('Alice', 22, 2, 40000, 1000),
('Bob', 28, 2, 60000, 2000),
('Charlie', 35, 3, 70000, 3000),
('David', 40, 3, 80000, 4000),
('Eve', 20, 4, 30000, 500),
('Frank', 45, 4, 90000, 6000);

INSERT IGNORE INTO Department (DepartName) VALUES
('Human Resources'),
('Finance'),
('Marketing'),
('Sales'),
('IT & Development'),
('Art and Design'),
('Customer Service'),
('Management'),
('Research and Development'),
('Quality Assurance'),
('Logistics');

INSERT IGNORE INTO Student (Name, Age, Grade) VALUES
('Roop', 32, 10),
('Adithya', 27, 9),
('Advik', 23, 8),
('Pratibha', 29, 9),
('Ashvitha', 24, 8),
('Dhruv', 30, 9),
('Dhariya Jaishwal', 35, 10),
('Chatur', 40, 11),
('Sowlaasya', 20, 7),
('Geetika', 45, 12),
('Pooja', 22, 8),
('Prachi', 28, 9),
('Sakshi', 35, 10),
('John Doe', 30, 9),
('Jane Doe', 25, 8),
('Alice', 22, 7),
('Bob', 28, 8),
('Charlie', 35, 9),
('David', 40, 10),
('Eve', 20, 6),
('Frank', 45, 11);

`

func sample(bot *telegram.BotAPI, message *telegram.Message) {
	if !accountCreateReminder(bot, message) {
		return
	}
	msg := telegram.NewMessage(message.Chat.ID, "")
	user := management.UserFromTelegram(message.From, &DB)
	user_db := user.GetDB(&DB)
	user_db.Connect()
	defer user_db.Disconnect()
	log.Println("Creating sample tables")
	user_db.UseDatabase(user.SQLDBName)

	queries := strings.Split(SampleQueries, ";")
	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}
		_, err := user_db.Conn.Exec(query)
		if err != nil {
			msg.Text = fmt.Sprintf("Error executing query: %s", err)
			log.Println(err)
		}
	}

	inserts := strings.Split(SampleInserts, ";")
	for _, insert := range inserts {
		insert = strings.TrimSpace(insert)
		if insert == "" {
			continue
		}
		_, err := user_db.Conn.Exec(insert)
		if err != nil {
			msg.Text = fmt.Sprintf("Error executing query: %s", err)
			log.Println(err)
		}
	}
	msg.Text = "Sample tables created successfully\n"
	msg.Text += "Use /sql to run queries\n"
	msg.Text += "Examples: \n`/sql SELECT * FROM Employee;`\n`/sql SELECT * FROM Student;`"
	msg.ParseMode = "MarkdownV2"
	bot.Send(msg)
}

func init() {
	Register(Command{
		Name:        "sample",
		Description: "Get a few sample tables \\(filled with sample data\\) to play with",
		Handler:     sample,
		Usage:       "/sample",
	})
}
