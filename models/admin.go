package models

type Admin struct {
	Id            int    `gorm:"primaryKey" json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Refresh_token string `json:"refresh_token"`
}

type AdminLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type AdminRegister struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
