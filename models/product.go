package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Images      []Image `json:"image_id"`
	Category_id uint    `json:"category_id" `
	Name        string  `json:"name" gorm:"unique;not null"`
	Price       int     `json:"price"`
	Quantity    int     `json:"quantity"`
	Specs       string  `json:"specs"`
}

type Category struct {
	gorm.Model
	// Product_id    []Product
	Category_name string `gorm:"unique;not null" json:"category_name" `
	Description   string `json:"description"`
	Status        bool   `json:"status"`
}

type ProductRespHome struct {
	ID          uint     `gorm:"primarykey" json:"id"`
	Images      []string `json:"image_id" gorm:"column:photo"`
	Category_id int      `json:"category_id" `
	Name        string   `json:"name"`
	Price       int      `json:"price"`
	// Quantity    int      `json:"quantity"`
	// Specs       string   `json:"specs"`
}
type ProductRespHomeDemo struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Images      string `json:"image_id" gorm:"column:photo"`
	Category_id int    `json:"category_id" `
	Name        string `json:"name"`
	Price       int    `json:"price"`
	// Quantity    int      `json:"quantity"`
	// Specs       string   `json:"specs"`
}
type ProductResp struct {
	ID          uint     `gorm:"primarykey" json:"id"`
	Images      []string `json:"image_id" gorm:"column:photo"`
	Category_id int      `json:"category_id" `
	Name        string   `json:"name"`
	Price       int      `json:"price"`
	Quantity    int      `json:"quantity"`
	Specs       string   `json:"specs"`
}
