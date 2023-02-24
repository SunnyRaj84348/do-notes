package models

import (
	"database/sql"

	"github.com/SunnyRaj84348/do-notes/initializers"
)

type Note struct {
	NoteTitle string `json:"noteTitle" binding:"required"`
	NoteBody  string `json:"noteBody"`
}

type Notes struct {
	NoteID    int    `json:"noteID"`
	NoteTitle string `json:"noteTitle"`
	NoteBody  string `json:"noteBody"`
}

func InsertNotes(userid int, noteTitle string, noteBody string) (int, error) {
	_, err := initializers.GetDB().Exec(`
		INSERT INTO notes(note_title, note_body, user_id) VALUES
		(?, ?, ?)
	`, noteTitle, noteBody, userid)

	if err != nil {
		return -1, err
	}

	row := initializers.GetDB().QueryRow(`SELECT LAST_INSERT_ID()`)

	var noteID int
	err = row.Scan(&noteID)

	return noteID, err
}

func GetNotes(userid int) (*sql.Rows, error) {
	rows, err := initializers.GetDB().Query(`
		SELECT note_id, note_title, note_body FROM notes
		WHERE user_id = ?
	`, userid)

	return rows, err
}

func UpdateNotes(userid int, noteID int, noteTitle string, noteBody string) error {
	row := initializers.GetDB().QueryRow(`SELECT user_id FROM notes WHERE note_id = ?`, noteID)
	var val int

	err := row.Scan(&val)
	if err == sql.ErrNoRows || userid != val {
		return sql.ErrNoRows
	}

	_, err = initializers.GetDB().Exec(`
		UPDATE notes SET note_title = ?, note_body = ?
		WHERE note_id = ? AND user_id = ?
	`, noteTitle, noteBody, noteID, userid)

	return err
}

func DeleteNotes(userid int, noteID int) error {
	row := initializers.GetDB().QueryRow(`SELECT user_id FROM notes WHERE note_id = ?`, noteID)
	var val int

	err := row.Scan(&val)
	if err == sql.ErrNoRows || userid != val {
		return sql.ErrNoRows
	}

	_, err = initializers.GetDB().Exec(`DELETE FROM notes WHERE note_id = ?`, noteID)
	return err
}
