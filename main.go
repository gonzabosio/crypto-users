package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gonzabosio/crypto-users/data"
	"github.com/gonzabosio/crypto-users/routes"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env data:", err)
	}

	data.Init()

	e := echo.New()

	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:x-apikey",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == os.Getenv("API_KEY"), nil
		},
	}))

	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "Crypto Users API") })

	e.GET("/users", routes.GetUserLogin)

	e.POST("/users", routes.PostUser)

	e.POST("/actions", routes.PostAction)

	e.GET("/actions", routes.GetAllUserActions)

	e.GET("/actions/:id", routes.GetActions)

	e.PATCH("/actions/:id", routes.PatchAction)

	e.DELETE("/actions/:id", routes.DeleteAction)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
