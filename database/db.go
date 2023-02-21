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

func Init(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS user(
			user_id INT PRIMARY KEY AUTO_INCREMENT,
			username VARCHAR(255) UNIQUE,
			password VARCHAR(255)
		)
	`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notes(
			note_id INT PRIMARY KEY AUTO_INCREMENT,
			note_title TEXT,
			note_body TEXT,
			user_id INT,
			FOREIGN KEY(user_id) REFERENCES user(user_id)
		)
	`)

	return err
}

func InsertUser(db *sql.DB, username string, password string) error {
	_, err := db.Exec(`INSERT INTO user(username, password) VALUES(?, ?)`, username, password)
	return err
}

func GetUser(db *sql.DB, username string) *sql.Row {
	row := db.QueryRow(`SELECT * FROM user WHERE username = ?`, username)
	return row
}

func InsertNotes(db *sql.DB, userid int, noteTitle string, noteBody string) (int, error) {
	_, err := db.Exec(`
		INSERT INTO notes(note_title, note_body, user_id) VALUES
		(?, ?, ?)
	`, noteTitle, noteBody, userid)

	if err != nil {
		return -1, err
	}

	row := db.QueryRow(`SELECT LAST_INSERT_ID()`)

	var noteID int
	err = row.Scan(&noteID)

	return noteID, err
}

func GetNotes(db *sql.DB, userid int) (*sql.Rows, error) {
	rows, err := db.Query(`
		SELECT note_id, note_title, note_body FROM notes
		WHERE user_id = ?
	`, userid)

	return rows, err
}

func UpdateNotes(db *sql.DB, userid int, noteID int, noteTitle string, noteBody string) error {
	row := db.QueryRow(`SELECT user_id FROM notes WHERE note_id = ?`, noteID)
	var val int

	err := row.Scan(&val)
	if err == sql.ErrNoRows || userid != val {
		return sql.ErrNoRows
	}

	_, err = db.Exec(`
		UPDATE notes SET note_title = ?, note_body = ?
		WHERE note_id = ? AND user_id = ?
	`, noteTitle, noteBody, noteID, userid)

	return err
}

func DeleteNotes(db *sql.DB, userid int, noteID int) error {
	row := db.QueryRow(`SELECT user_id FROM notes WHERE note_id = ?`, noteID)
	var val int

	err := row.Scan(&val)
	if err == sql.ErrNoRows || userid != val {
		return sql.ErrNoRows
	}

	_, err = db.Exec(`DELETE FROM notes WHERE note_id = ?`, noteID)
	return err
}
