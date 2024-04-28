package helpers

import (
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), viper.GetInt("BCRYPT_SALT"))
	if err != nil {
		return "", err
	}
	return string(hashedPass), nil
}

func ComparePassword(hashPassword string, password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		return false, err
	}
	return true, nil
}
