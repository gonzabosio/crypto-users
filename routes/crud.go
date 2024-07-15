package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAction(c echo.Context) error {
	return c.String(http.StatusOK, "List of users/actions")
}

func PostAction(c echo.Context) error {
	return c.String(http.StatusOK, "Add an action")
}

func PatchAction(c echo.Context) error {
	return c.String(http.StatusOK, "Edit an action")
}

func DeleteAction(c echo.Context) error {
	return c.String(http.StatusOK, "Delete an action")
}
