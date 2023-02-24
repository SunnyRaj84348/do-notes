package initializers

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func ConnectToDB() error {
	var err error

	db, err = sql.Open("mysql", os.Getenv("DB_URL"))
	if err != nil {
		return err
	}

	err = db.Ping()
	return err
}

func InitDB() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS user(
			user_id INT PRIMARY KEY AUTO_INCREMENT,
			username VARCHAR(255) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL
		)
	`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notes(
			note_id INT PRIMARY KEY AUTO_INCREMENT,
			note_title TEXT NOT NULL,
			note_body TEXT,
			user_id INT NOT NULL,
			FOREIGN KEY(user_id) REFERENCES user(user_id)
		)
	`)

	return err
}

func GetDB() *sql.DB {
	return db
}
