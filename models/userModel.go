package models

import (
	"time"
)

type Credential struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	UserID    string `gorm:"primaryKey; type:uuid; default:gen_random_uuid()"`
	Email     string `gorm:"type:citext; unique; not null" json:"email" binding:"required"`
	Username  string `gorm:"unique; not null" json:"username" binding:"required"`
	Password  string `gorm:"not null" json:"password" binding:"required"`
	Verified  bool
	CreatedAt time.Time
}

type EmailAuth struct {
	Email     string    `gorm:"primaryKey; type:citext" json:"email" binding:"required"`
	Code      string    `gorm:"type:char(4) not null " json:"code" binding:"required"`
	ExpiresAt time.Time `gorm:"not null"`
}

func InsertUser(user User) error {
	user.Verified = false
	tx := db.Create(&user)

	return tx.Error
}

func GetUser(username string) (User, error) {
	user := User{}
	tx := db.First(&user, "username = ? OR email = ?", username, username)

	return user, tx.Error
}

func InsertEmailAuth(email string, code string) error {
	emailAuth := EmailAuth{email, code, time.Now().Add(10 * time.Minute)}
	tx := db.Create(&emailAuth)

	return tx.Error
}

func GetEmailAuth(emailAuth EmailAuth) (EmailAuth, error) {
	tx := db.First(&emailAuth, "code = ?", emailAuth.Code)
	return emailAuth, tx.Error
}

func DeleteEmailAuth(emailAuth EmailAuth) error {
	tx := db.Delete(&EmailAuth{}, "email = ?", emailAuth.Email)
	return tx.Error
}

func SetUserVerified(email string) error {
	user, err := GetUser(email)
	if err != nil {
		return err
	}

	user.Verified = true

	tx := db.Save(&user)
	return tx.Error
}
