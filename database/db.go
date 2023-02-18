package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Connect(conn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	return db, err
}

func InsertUser(db *sql.DB, username string, password string) error {
	_, err := db.Exec(`INSERT INTO user(username, password) VALUES(?, ?)`, username, password)
	if err != nil {
		return err
	}

	return nil
}

func GetUser(db *sql.DB, username string) *sql.Row {
	row := db.QueryRow(`SELECT * FROM user WHERE username = ?`, username)
	return row
}

func InsertNotes(db *sql.DB, username string, noteTitle string, noteBody string) error {
	_, err := db.Exec(`
		INSERT INTO notes(note_title, note_body, user_id) VALUES
		(?, ?, (SELECT user_id FROM user WHERE username = ?))
	`, noteTitle, noteBody, username)

	return err
}
