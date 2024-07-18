package routes

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gonzabosio/crypto-users/data"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func UserLogin(c echo.Context) error {
	var user data.User
	err := c.Bind(&user)
	if err != nil {
		return c.String(http.StatusInternalServerError, "body request binding failed")
	}
	var dbUser data.User
	result := data.DB.Where("username = ?", user.Username).First(&dbUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.String(http.StatusNotFound, "User not found")
		}
		return c.String(http.StatusInternalServerError, "get user data failed")
	}
	if err = data.VerifyPassword(dbUser.Password, user.Password); err != nil {
		return c.String(http.StatusUnauthorized, "Incorrect password")
	}
	userData := data.UserData{
		ID:       fmt.Sprintf("%v", dbUser.ID),
		Username: dbUser.Username,
	}
	return c.JSON(http.StatusOK, userData)
}

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
	if err := result.Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return c.String(http.StatusConflict, "Username already exists")
		}
		return c.String(http.StatusInternalServerError, "record creation failed")
	}
	userData := data.UserData{
		ID:       fmt.Sprintf("%v", user.ID),
		Username: user.Username,
	}
	return c.JSON(http.StatusOK, userData)
}

func PostAction(c echo.Context) error {
	var act data.ActivityPostAdapted
	err := c.Bind(&act)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request: "+err.Error())
	}
	id, err := strconv.Atoi(act.UserID)
	if err != nil {
		log.Fatal("convert id failed" + err.Error())
	}
	activity := data.Activity{
		Action:       act.Action,
		CryptoCode:   act.CryptoCode,
		CryptoAmount: act.CryptoAmount,
		Money:        act.Money,
		CreatedAt:    act.CreatedAt,
		UserID:       uint(id),
	}
	result := data.DB.Create(&activity)
	if result.Error != nil {
		return c.String(http.StatusInternalServerError, "activity creation failed: "+result.Error.Error())
	}
	return c.String(http.StatusOK, "activity created successfully")
}

func GetActions(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.String(http.StatusInternalServerError, "id to int failed")
	}
	var dbActivities []data.Activity
	result := data.DB.Where("\"Activity\".user_id = ?", id).Find(&dbActivities)
	if err = result.Error; err != nil {
		return fmt.Errorf("actions of user get method failed: " + err.Error())
	}
	var activities []data.ActivityGetAdapted
	for _, activity := range dbActivities {
		activityGetMethod := data.ActivityGetAdapted{
			ID:           fmt.Sprintf("%v", activity.ID),
			Action:       activity.Action,
			CryptoCode:   activity.CryptoCode,
			CryptoAmount: activity.CryptoAmount,
			Money:        activity.Money,
			CreatedAt:    activity.CreatedAt,
			UserID:       fmt.Sprintf("%v", activity.UserID),
		}
		activities = append(activities, activityGetMethod)
	}
	return c.JSON(http.StatusOK, activities)
}

func GetAllUserActions(c echo.Context) error {
	var users []data.User
	result := data.DB.Table("\"User\"").
		Select("\"User\".id, \"User\".username").
		Joins("JOIN \"Activity\" ON \"User\".id=\"Activity\".user_id").Group("\"User\".id").Find(&users)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.String(http.StatusNotFound, "actions not found")
		}
		return c.String(http.StatusInternalServerError, "failed to retrieve actions")
	}
	for i := range users {
		var activities []data.Activity
		data.DB.Where("user_id = ?", users[i].ID).Find(&activities)
		users[i].Activity = activities
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
			return c.String(http.StatusNotFound, "action not found")
		}
		return c.String(http.StatusInternalServerError, "Failed to retrieve action")
	}

	var updateAction data.PatchActivity
	if err = c.Bind(&updateAction); err != nil {
		return c.String(http.StatusBadRequest, "failed to bind request")
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
