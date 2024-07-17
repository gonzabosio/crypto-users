package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gonzabosio/crypto-users/data"
	"github.com/gonzabosio/crypto-users/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	data.Init()

	e := echo.New()

	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 10, Burst: 10, ExpiresIn: 3 * time.Minute},
		),
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}
	e.Use(middleware.RateLimiterWithConfig(config))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{os.Getenv("HOST"), "https://crypto-bull-gb.vercel.app"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "x-apikey"},
		AllowMethods: []string{echo.GET, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:x-apikey",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == os.Getenv("API_KEY"), nil
		},
	}))

	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "Crypto Users API") })

	e.POST("/register", routes.PostUser)

	e.POST("/login", routes.UserLogin)

	e.POST("/actions", routes.PostAction)

	e.GET("/actions", routes.GetAllUserActions)

	e.GET("/actions/:id", routes.GetActions)

	e.PATCH("/actions/:id", routes.PatchAction)

	e.DELETE("/actions/:id", routes.DeleteAction)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
