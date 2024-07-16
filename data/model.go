package data

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type User struct {
	ID       uint       `json:"user_id" gorm:"primaryKey"`
	Username string     `json:"username" gorm:"unique"`
	Password string     `json:"password"`
	Activity []Activity `json:"activity" gorm:"foreignKey:UserID"`
}

type Activity struct {
	ID           uint       `json:"activity_id" gorm:"primaryKey"`
	Action       string     `json:"action"`
	CryptoCode   string     `json:"crypto_code"`
	CryptoAmount float32    `json:"crypto_amount"`
	Money        float32    `json:"money"`
	CreatedAt    CustomTime `json:"performed_at"`
	UserID       uint       `json:"user_id"`
}

func (User) TableName() string {
	return "User"
}

func (Activity) TableName() string {
	return "Activity"
}

type CustomTime struct {
	time.Time
}

func (ct CustomTime) Value() (driver.Value, error) {
	return ct.Time, nil
}

func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		*ct = CustomTime{Time: time.Time{}}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*ct = CustomTime{Time: v}
		return nil
	case string:
		parsedTime, err := time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			return err
		}
		*ct = CustomTime{Time: parsedTime}
		return nil
	default:
		return fmt.Errorf("cannot scan %T into CustomTime", value)
	}
}

func (ct *CustomTime) ParseDate(b []byte) error {
	str := strings.Trim(string(b), `"`)
	loc, err := time.LoadLocation("America/Buenos_Aires")
	if err != nil {
		return err
	}
	parsedTime, err := time.ParseInLocation("2006-01-02 15:04:05", str, loc)
	if err != nil {
		return err
	}
	ct.Time = parsedTime
	return nil
}

type UpdateActivity struct {
	Action       *string  `json:"action,omitempty"`
	CryptoCode   *string  `json:"crypto_code,omitempty"`
	CryptoAmount *float32 `json:"crypto_amount,omitempty"`
	Money        *float32 `json:"money,omitempty"`
}
