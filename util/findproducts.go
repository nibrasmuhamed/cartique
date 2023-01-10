package util

import (
	"fmt"

	"github.com/nibrasmuhamed/cartique/database"
	"github.com/nibrasmuhamed/cartique/models"
)

func FindProducts(a []int) []models.ProductRespHome {
	db := database.OpenDataBase()
	defer database.CloseDatabase(db)
	var p []models.ProductRespHome
	for _, val := range a {
		r, err := db.Query("SELECT id, name, price, category_id from products where products.deleted_at is null AND id=?", val)
		if err != nil {
			fmt.Println("error is :", err)
		}
		defer r.Close()
		for r.Next() {
			a := models.ProductRespHome{}
			err = r.Scan(&a.ID, &a.Name, &a.Price, &a.Category_id)
			if err != nil {
				fmt.Println("error while scanning ", err)
			}
			i, err := db.Query("SELECT photo FROM images WHERE images.product_id=?", a.ID)
			if err != nil {
				fmt.Println("2nd errror is:", err)
			}
			for i.Next() {
				var x string
				err = i.Scan(&x)
				if err != nil {
					fmt.Println(err)
				}
				a.Images = append(a.Images, x)
			}
			p = append(p, a)
		}
	}
	return p
}
