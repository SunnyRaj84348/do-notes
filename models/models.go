package models

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func ConnectToDB() error {
	var err error

	db, err = gorm.Open(postgres.Open(os.Getenv("DB_URL")), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	return err
}

func InitDB() error {
	if !db.Migrator().HasTable("user") || !db.Migrator().HasTable("notes") {
		err := db.AutoMigrate(User{}, Notes{})
		return err
	}

	return nil
}
