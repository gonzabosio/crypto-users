package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gonzabosio/crypto-users/data"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env data:", err)
	}

	dsn := os.Getenv("CONN_STR")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	} else {
		log.Printf("successful database connection")
	}

	user := new(data.User)

	err = db.AutoMigrate(&data.User{}, &data.Activity{})
	if err != nil {
		log.Fatal("auto migration failed")
	}
	db.Raw("SELECT username FROM \"User\"").Scan(&user)

	e := echo.New()

	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "Hello world") })

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
