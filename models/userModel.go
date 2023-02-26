package models

type Credential struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	UserID   uint32 `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
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
