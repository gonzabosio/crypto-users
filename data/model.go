package data

import "time"

type User struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	CryptoCode   float32   `json:"crypto_code"`
	CryptoAmount float32   `json:"crypto_amount"`
	Money        float32   `json:"money"`
	Action       string    `json:"action"`
	Datetime     time.Time `json:"datetime"`
}

func (User) TableName() string {
	return "User"
}
