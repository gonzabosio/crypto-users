package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gonzabosio/crypto-users/data"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func PostUser(c echo.Context) error {
	var user data.User
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	pw, err := data.HashPassword(user.Password)
	if err != nil {
		return c.String(http.StatusExpectationFailed, "hash password method failed")
	}
	user.Password = pw
	result := data.DB.Create(&user)
	if result.Error != nil {
		return c.String(http.StatusInternalServerError, "record creation failed")
	}
	response := fmt.Sprintf("new user added id: %v | username: %v", user.ID, user.Username)
	return c.String(http.StatusOK, response)
}

func PostAction(c echo.Context) error {
	var act data.Activity
	err := c.Bind(&act)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request: "+err.Error())
	}
	result := data.DB.Create(&act)
	if result.Error != nil {
		return c.String(http.StatusInternalServerError, "activity creation failed")
	}
	return c.String(http.StatusOK, "activity created successfully")
}

func GetActions(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusInternalServerError, "id to int failed")
	}
	var activities []data.Activity
	result := data.DB.Where("\"Activity\".user_id = ?", id).Find(&activities)
	if err = result.Error; err != nil {
		return fmt.Errorf("actions of user get method failed: " + err.Error())
	}
	return c.JSON(http.StatusOK, activities)
}

func GetAllUserActions(c echo.Context) error {
	var users data.User
	result := data.DB.Preload("Activity").
		Select("\"User\".id, \"User\".username").
		Joins("JOIN \"Activity\" ON \"User\".id=\"Activity\".user_id").Find(&users)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.String(http.StatusNotFound, "Actions not found")
		}
		return c.String(http.StatusInternalServerError, "Failed to retrieve actions")
	}

	return c.JSON(http.StatusOK, users)
}

func PatchAction(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusInternalServerError, "id to int failed")
	}

	var act data.Activity
	result := data.DB.First(&act, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.String(http.StatusNotFound, "Action not found")
		}
		return c.String(http.StatusInternalServerError, "Failed to retrieve action")
	}

	var updateAction data.UpdateActivity
	if err = c.Bind(&updateAction); err != nil {
		return c.String(http.StatusBadRequest, "Failed to bind request")
	}
	if updateAction.Action != nil {
		act.Action = *updateAction.Action
	}
	if updateAction.CryptoAmount != nil {
		act.CryptoAmount = *updateAction.CryptoAmount
	}
	if updateAction.CryptoCode != nil {
		act.CryptoCode = *updateAction.CryptoCode
	}
	if updateAction.Money != nil {
		act.Money = *updateAction.Money
	}
	result = data.DB.Save(&act)
	if result.Error != nil {
		return c.String(http.StatusInternalServerError, "failed to update activity")
	}

	return c.JSON(http.StatusOK, act)
}

func DeleteAction(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusInternalServerError, "id to int failed")
	}
	result := data.DB.Delete(&data.Activity{}, id)
	if result.Error != nil {
		return c.String(http.StatusNotFound, "action id does not exist")
	}
	return c.String(http.StatusOK, "action deleted")
}
