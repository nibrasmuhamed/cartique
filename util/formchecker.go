package util

import "github.com/nibrasmuhamed/cartique/models"

func CheckFiles(p models.Product) bool {
	if p.Category_id == 0 || p.Name == "" || p.Price == 0 || p.Quantity == 0 || p.Specs == "" {
		return false
	}
	return true
}
