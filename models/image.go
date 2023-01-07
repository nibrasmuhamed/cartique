package models

import "gorm.io/gorm"

type Image struct {
	gorm.Model
	Product_id uint   `json:"product_id"`
	Photo      string `json:"image_link"`
}

type ImageRespose struct {
	Photo []string `json:"image_link"`
}
