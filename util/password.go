package util

import (
	"log"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func Password(pass string) bool {
	var (
		upp, low, num, sym bool
		tot                uint8
	)

	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			upp = true
			tot++
		case unicode.IsLower(char):
			low = true
			tot++
		case unicode.IsNumber(char):
			num = true
			tot++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
			tot++
		default:
			return false
		}
	}

	if !upp || !low || !num || !sym || tot < 8 {
		return false
	}

	return true
}

func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error while hashing password: ", err)
	}
	return string(hash)
}

func VerifyPassword(password string, dbpassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dbpassword), []byte(password))
	if err != nil {
		log.Println("password do not match")
		return false
	}
	return true
}
