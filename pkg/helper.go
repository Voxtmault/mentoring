package pkg

import (
	"golang.org/x/crypto/bcrypt"
)

func CompareHashBcrypt(hashedInput string, input string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedInput), []byte(input))
}

func HashBcrypt(input string) (string, error) {
	hashedInput, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedInput), nil
}
