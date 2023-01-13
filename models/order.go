package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID    uint `gorm:"not null" json:"user_id"`
	User      User
	ProductID uint `json:"product_id"`
	Product   Product
	AddressID uint `json:"address_id"`
	Address   Address
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price"`
	Status    string `json:"status"`
}

type OrderResp struct {
	ID        uint   `json:"id"`
	ProductID uint   `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price"`
	Status    string `json:"status"`
}

type OrderRespAdmin struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	ProductID uint   `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Price     int    `json:"price"`
	Status    string `json:"status"`
}
