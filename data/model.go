package data

import (
	"time"
)

type User struct {
	ID       uint       `json:"user_id" gorm:"primaryKey"`
	Username string     `json:"username" gorm:"unique"`
	Password string     `json:"password"`
	Activity []Activity `json:"activity" gorm:"foreignKey:UserID"`
}

type Activity struct {
	ID           uint      `json:"act_id" gorm:"primaryKey"`
	Action       string    `json:"action"`
	CryptoCode   string    `json:"crypto_code"`
	CryptoAmount float32   `json:"crypto_amount"`
	Money        float32   `json:"money"`
	CreatedAt    time.Time `json:"made_at"`
	UpdatedAt    time.Time `json:"modified_at"`
	UserID       uint      `json:"user_id"`
}

func (User) TableName() string {
	return "User"
}

func (Activity) TableName() string {
	return "Activity"
}
