package data

import (
	"time"
)

type User struct {
	ID       int64      `json:"user_id" gorm:"primaryKey"`
	Username string     `json:"username" gorm:"unique"`
	Password string     `json:"password"`
	Activity []Activity `json:"activity" gorm:"foreignKey:UserID"`
}

type Activity struct {
	ID           int64     `json:"activity_id" gorm:"primaryKey"`
	Action       string    `json:"action"`
	CryptoCode   string    `json:"crypto_code"`
	CryptoAmount float32   `json:"crypto_amount"`
	Money        float32   `json:"money"`
	CreatedAt    time.Time `json:"performed_at"`
	UserID       uint      `json:"user_id"`
}

func (User) TableName() string {
	return "User"
}

func (Activity) TableName() string {
	return "Activity"
}

type ActivityGetAdapted struct {
	ID           string    `json:"activity_id" gorm:"primaryKey"`
	Action       string    `json:"action"`
	CryptoCode   string    `json:"crypto_code"`
	CryptoAmount float32   `json:"crypto_amount"`
	Money        float32   `json:"money"`
	CreatedAt    time.Time `json:"performed_at"`
	UserID       string    `json:"user_id"`
}

type ActivityPostAdapted struct {
	Action       string    `json:"action"`
	CryptoCode   string    `json:"crypto_code"`
	CryptoAmount float32   `json:"crypto_amount"`
	Money        float32   `json:"money"`
	CreatedAt    time.Time `json:"performed_at"`
	UserID       string    `json:"user_id"`
}

type PatchActivity struct {
	Action       *string  `json:"action,omitempty"`
	CryptoCode   *string  `json:"crypto_code,omitempty"`
	CryptoAmount *float32 `json:"crypto_amount,omitempty"`
	Money        *float32 `json:"money,omitempty"`
}

type UserData struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
}
