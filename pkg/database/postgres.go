package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewPostgresDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to the database!")
	return db, nil
}

func CreateUserTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        phone_number VARCHAR(20) UNIQUE NOT NULL,
        registered_at TIMESTAMPTZ NOT NULL
    );`
	_, err := db.Exec(query)
	return err
}
