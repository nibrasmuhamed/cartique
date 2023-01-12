package models

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	// foreign key of user
	UserID   uint
	User     User   // `gorm:"foreignKey:ID"`
	Home     string `json:"home"`
	City     string `json:"city"`
	District string `json:"district"`
	State    string `json:"state"`
	Pin      int    `json:"pin"`
}

type AddressResp struct {
	ID       uint   `json:"id"`
	UserID   uint   `json:"user_id"`
	Home     string `json:"home"`
	City     string `json:"city"`
	District string `json:"district"`
	State    string `json:"state"`
	Pin      int    `json:"pin"`
}
