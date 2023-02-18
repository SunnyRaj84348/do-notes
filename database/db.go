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
	return err
}

func GetUser(db *sql.DB, username string) *sql.Row {
	row := db.QueryRow(`SELECT * FROM user WHERE username = ?`, username)
	return row
}

func InsertNotes(db *sql.DB, userid int, noteTitle string, noteBody string) error {
	_, err := db.Exec(`
		INSERT INTO notes(note_title, note_body, user_id) VALUES
		(?, ?, ?)
	`, noteTitle, noteBody, userid)

	return err
}

func GetNotes(db *sql.DB, userid int) (*sql.Rows, error) {
	rows, err := db.Query(`
		SELECT note_id, note_title, note_body FROM notes
		WHERE user_id = ?
	`, userid)

	return rows, err
}
