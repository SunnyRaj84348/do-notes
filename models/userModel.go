package models

import (
	"database/sql"

	"github.com/SunnyRaj84348/do-notes/initializers"
)

type Credential struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	UserID   int
	Username string
	Password string
}

func InsertUser(username string, password string) error {
	_, err := initializers.GetDB().Exec(`INSERT INTO user(username, password) VALUES(?, ?)`, username, password)
	return err
}

func GetUser(username string) *sql.Row {
	row := initializers.GetDB().QueryRow(`SELECT * FROM user WHERE username = ?`, username)
	return row
}
