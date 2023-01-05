package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	TokenExpired = errors.New("token is expired")
	UnknownError = errors.New("internal server error")
)

func CheckToken(t string) (bool, error) {
	key := []byte(os.Getenv("KEY"))
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return key, nil
	})
	if err != nil {
		log.Println(err)
		return false, TokenExpired
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["role"] == "user" {
			return true, nil
		}
	}
	return false, UnknownError
}

func GetidfromToken(t string) float64 {
	key := []byte(os.Getenv("KEY"))
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return key, nil
	})
	if err != nil {
		log.Println(err)
		return 0
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["role"] == "user" {
			return claims["id"].(float64)
		}
	}
	return 0
}

func GenerateJWT(id int, role string) (string, error) {
	var mySigningKey = []byte(os.Getenv("KEY"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = id
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Printf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}
