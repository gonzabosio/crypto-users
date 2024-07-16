package data

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := os.Getenv("CONN_STR")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	} else {
		log.Printf("successful database connection")
	}
	err = DB.AutoMigrate(&User{}, &Activity{})
	if err != nil {
		log.Fatal("auto migration failed")
	}
}
