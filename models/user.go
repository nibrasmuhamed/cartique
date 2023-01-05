package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username      string `json:"username"`
	Email         string `json:"email" gorm:"index"`
	Phone         string `json:"phone"`
	Password      string `json:"password"`
	Verified      bool   `json:"verified"`
	Otp           int    `json:"otp"`
	Refresh_token string `json:"refresh_token"`
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserRegister struct {
	Name     string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
