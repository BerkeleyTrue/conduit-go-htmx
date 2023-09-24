package password

import (
	"errors"

	"github.com/berkeleytrue/conduit/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type (
	Password       string
	HashedPassword string
)

var (
	passwordTooShort           error = errors.New("Password must be at least 8 characters long")
	passwordTooLong            error = errors.New("Password must be at most 64 characters long")
	passwordMissingLowercase   error = errors.New("Password must contain at least one lowercase character")
	passwordMissingUppercase   error = errors.New("Password must contain at least one uppercase character")
	passwordMissingNumber      error = errors.New("Password must contain at least one number")
	passwordMissingSpecialChar error = errors.New("Password must contain at least one special character")
)

func New(raw string) (Password, error) {
	passwordLength := len(raw)
	pass := []rune(raw)

	if passwordLength < 8 {
		return "", passwordTooShort

	} else if passwordLength > 64 {
		return "", passwordTooLong

	} else if !utils.IsAnyLowercase(pass) {
		return "", passwordMissingLowercase

	} else if !utils.IsAnyUppercase(pass) {
		return "", passwordMissingUppercase

	} else if !utils.IsAnyDigit(pass) {
		return "", passwordMissingNumber

	} else if !utils.IsAnySymbol(pass) {
		return "", passwordMissingSpecialChar
	}

	return Password(pass), nil
}

func HashPassword(pass Password) (HashedPassword, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return HashedPassword(hashed), nil
}

func CompareHashAndPassword(hashed HashedPassword, pass Password) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pass))
}
