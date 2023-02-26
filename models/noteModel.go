package models

type Note struct {
	NoteTitle string `json:"noteTitle" binding:"required"`
	NoteBody  string `json:"noteBody"`
}

type Notes struct {
	NoteID    uint32 `gorm:"primaryKey" json:"noteID"`
	NoteTitle string `gorm:"not null" json:"noteTitle"`
	NoteBody  string `json:"noteBody"`
	UserID    uint32 `gorm:"not null" json:"-"`
	User      User   `json:"-"`
}

func InsertNotes(userid uint32, noteTitle string, noteBody string) (Notes, error) {
	notes := Notes{NoteTitle: noteTitle, NoteBody: noteBody, UserID: userid}
	tx := db.Create(&notes)

	return notes, tx.Error
}

func GetNotes(userid uint32) ([]Notes, error) {
	notes := []Notes{}
	tx := db.Find(&notes, "user_id = ?", userid)

	return notes, tx.Error
}

func UpdateNotes(userid uint32, noteID uint32, noteTitle string, noteBody string) (Notes, error) {
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

func DeleteNotes(userid uint32, noteID uint32) error {
	notes := Notes{NoteID: noteID}

	tx := db.First(&notes, "user_id = ?", userid)
	if tx.Error != nil {
		return tx.Error
	}

	tx = db.Delete(&notes)
	return tx.Error
}
