package infra

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("DB not reachable:", err)
	}

	db.SetMaxOpenConns(25)                 // Max number of open connections at any time
	db.SetMaxIdleConns(10)                 // Max idle connections kept ready
	db.SetConnMaxLifetime(time.Minute * 5) // How long a connection is reused before closing

	schema := `
		CREATE TABLE IF NOT EXISTS groups (id SERIAL PRIMARY KEY, name TEXT NOT NULL);
		CREATE TABLE IF NOT EXISTS users (id SERIAL PRIMARY KEY, name TEXT, password TEXT NOT NULL);
		CREATE TABLE IF NOT EXISTS expenses (
			id SERIAL PRIMARY KEY, 
			group_id INT REFERENCES groups(id), -- Added this!
			paid_by INT REFERENCES users(id), 
			amount NUMERIC
		);
		CREATE TABLE IF NOT EXISTS participants (expense_id INT REFERENCES expenses(id), user_id INT REFERENCES users(id));
		CREATE TABLE IF NOT EXISTS group_members (
			group_id INT REFERENCES groups(id),
			user_id INT REFERENCES users(id),
			PRIMARY KEY (group_id, user_id)
		);
	`

	_, err = db.Exec(schema)
	if err != nil {
		log.Fatal("Schema error:", err)
	}

	indexes := `
		CREATE INDEX IF NOT EXISTS idx_expenses_paid_by ON expenses(paid_by);
		CREATE INDEX IF NOT EXISTS idx_participants_expense_id ON participants(expense_id);
		CREATE INDEX IF NOT EXISTS idx_group_members_user_id ON group_members(user_id);
	`

	_, err = db.Exec(indexes)
	if err != nil {
		log.Fatal("Index creation error:", err)
	}

	fmt.Println("Connected to PostgresSQL")

	return db
}
