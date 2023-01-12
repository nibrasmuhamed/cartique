package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	ProductID uint
	Product   Product //`gorm:"foreignKey:ID"`
	Quantity  int
	UserID    uint
	User      User
	// Total     int
}

type Wishlist struct {
	gorm.Model
	ProductID uint
	Product   Product //`gorm:"foreignKey:ID"`
	Quantity  int
	UserID    uint
	// User      User `gorm:"foreignKey:ID"`
	// Total     int
}
