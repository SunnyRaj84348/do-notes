package models

import "time"

type Credential struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type User struct {
	UserID    string `gorm:"primaryKey; type:uuid; default:gen_random_uuid()"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	CreatedAt time.Time
}

func InsertUser(username string, password string) error {
	user := User{Username: username, Password: password}
	tx := db.Create(&user)

	return tx.Error
}

func GetUser(username string) (User, error) {
	user := User{}
	tx := db.First(&user, "username = ?", username)

	return user, tx.Error
}
