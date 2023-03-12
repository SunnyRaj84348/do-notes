package models

import (
	"html/template"
	"time"

	"gorm.io/gorm"
)

type Note struct {
	NoteTitle string        `json:"noteTitle" binding:"required"`
	NoteBody  template.HTML `json:"noteBody"`
}

type Notes struct {
	NoteID    string         `gorm:"primaryKey; type:uuid; default:gen_random_uuid()" json:"noteID"`
	NoteTitle string         `gorm:"not null" json:"noteTitle"`
	NoteBody  template.HTML  `json:"noteBody"`
	UserID    string         `gorm:"not null" json:"-"`
	User      User           `json:"-"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func InsertNotes(userid string, noteTitle string, noteBody template.HTML) (Notes, error) {
	notes := Notes{NoteTitle: noteTitle, NoteBody: noteBody, UserID: userid}
	tx := db.Create(&notes)

	return notes, tx.Error
}

func GetNotes(userid string) ([]Notes, error) {
	notes := []Notes{}
	tx := db.Find(&notes, "user_id = ?", userid)

	return notes, tx.Error
}

func UpdateNotes(userid string, noteID string, noteTitle string, noteBody template.HTML) (Notes, error) {
	notes := Notes{NoteID: noteID}

	tx := db.First(&notes, "user_id = ?", userid)
	if tx.Error != nil {
		return notes, tx.Error
	}

	notes.NoteTitle = noteTitle
	notes.NoteBody = noteBody

	tx = db.Save(&notes)
	return notes, tx.Error
}

func DeleteNotes(userid string, noteID string) error {
	notes := Notes{NoteID: noteID}

	tx := db.First(&notes, "user_id = ?", userid)
	if tx.Error != nil {
		return tx.Error
	}

	tx = db.Delete(&notes)
	return tx.Error
}
