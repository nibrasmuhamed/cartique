package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

var (
	ErrTokenExpired = errors.New("token is expired")
	ErrUnknownError = errors.New("internal server error")
)

func CheckToken(t string) (bool, error) {
	key := []byte(os.Getenv("KEY"))
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing")
		}
		return key, nil
	})
	if err != nil {
		log.Println(err)
		return false, ErrTokenExpired
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["role"] == "user" {
			return true, nil
		}
	}
	return false, ErrUnknownError
}

func CheckTokenAdmin(t string) (bool, error) {
	key := []byte(os.Getenv("KEY"))
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing")
		}
		return key, nil
	})
	if err != nil {
		log.Println(err)
		return false, ErrTokenExpired
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["role"] == "admin" {
			return true, nil
		}
	}
	return false, ErrUnknownError
}

func GetidfromToken(t string) float64 {
	key := []byte(os.Getenv("KEY"))
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing")
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

func GenerateJWT(id int, role string, duration int) (string, error) {
	var mySigningKey = []byte(os.Getenv("KEY"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = id
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(duration)).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Printf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func GetAuthRef(c *fiber.Ctx) (string, string, error) {
	t := c.Get("Authorization")
	r := c.Get("refresh_token")
	if t == "" || r == "" {
		return "", "", errors.New("not found r and t")
	}
	return t[7:], r, nil
}
