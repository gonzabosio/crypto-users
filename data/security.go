package data

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "hash error", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword string, userEnteredPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userEnteredPassword))
}
