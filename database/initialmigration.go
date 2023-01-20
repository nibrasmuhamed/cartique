package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nibrasmuhamed/cartique/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// var DB *gorm.DB

func InitDB() *gorm.DB {

	DB, err := gorm.Open(mysql.Open(os.Getenv("DB")), &gorm.Config{})
	if err != nil {
		log.Println("error in connecting database : ", err)
	}
	DB.AutoMigrate(models.User{})
	DB.AutoMigrate(models.Admin{})
	DB.AutoMigrate(models.Product{})
	DB.AutoMigrate(models.Category{})
	DB.AutoMigrate(models.Image{})
	DB.AutoMigrate(models.Cart{})
	DB.AutoMigrate(models.Address{})
	DB.AutoMigrate(models.Wishlist{})
	DB.AutoMigrate(models.Order{})

	return DB
}

func CloseDb(db *gorm.DB) {
	dbInstance, _ := db.DB()
	_ = dbInstance.Close()
}

func OpenDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open(os.Getenv("DB")), &gorm.Config{})
	if err != nil {
		log.Println("error in connecting database : ", err)
	}
	return db
}

func OpenDataBase() *sql.DB {
	db, err := sql.Open("mysql", os.Getenv("DB"))
	if err != nil {
		fmt.Println("error is :", err)
	}
	return db
}

func CloseDatabase(d *sql.DB) {
	d.Close()
}

// Jan 20 06:04:24 ip-172-31-41-198 cartique[8244]: 2023/01/20 06:04:24 /home/ubuntu/cartique/database/initialmigration.go:26 Error 3780 (HY000): Referencing column 'category_id' and referenced column 'id' in foreign
//  key constraint 'fk_categories_product_id' are incompatible.
// Jan 20 06:04:24 ip-172-31-41-198 cartique[8244]: [8.822ms] [rows:0] ALTER TABLE `products` ADD CONSTRAINT `fk_categories_product_id` FOREIGN KEY (`category_id`) REFERENCES `categories`(`id`)
// Jan 20 06:04:24 ip-172-31-41-198 cartique[8244]:
// Jan 20 06:04:24 ip-172-31-41-198 cartique[8244]: 2023/01/20 06:04:24 /home/ubuntu/cartique/database/initialmigration.go:27 Error 3780 (HY000): Referencing column 'category_id' and referenced column 'id' in foreign
//  key constraint 'fk_categories_product_id' are incompatible.
// Jan 20 06:04:24 ip-172-31-41-198 cartique[8244]: [0.645ms] [rows:0] ALTER TABLE `products` ADD CONSTRAINT `fk_categories_product_id` FOREIGN KEY (`category_id`) REFERENCES `categories`(`id`)
// Jan 20 06:04:24 ip-172-31-41-198 cartique[8244]:
// Jan 20 06:04:24 ip-172-31-41-198 cartique[8244]: 2023/01/20 06:04:24 /home/ubuntu/cartique/database/initialmigration.go:29 Error 3780 (HY000): Referencing column 'category_id' and referenced column 'id' in foreign
//  key constraint 'fk_categories_product_id' are incompatible.
// Jan 20 06:04:24 ip-172-31-41-198 cartique[8244]: [0.554ms] [rows:0] ALTER TABLE `products` ADD CONSTRAINT `fk_categories_product_id` FOREIGN KEY (`category_id`) REFERENCES `categories`(`id`)
// Jan 20 06:04:24 ip-172-31-41-198 cartique[8244]:
// Jan 20 06:04:24 ip-172-31-41-198 cartique[8244]: 2023/01/20 06:04:24 /home/ubuntu/cartique/database/initialmigration.go:30 Error 3780 (HY000): Referencing column 'category_id' and referenced column 'id' in foreign
//  key constraint 'fk_categories_product_id' are incompatible.
// Jan 20 06:04:24 ip-172-31-41-198 cartique[8244]: [0.512ms] [rows:0] ALTER TABLE `products` ADD CONSTRAINT `fk_categories_product_id` FOREIGN KEY (`category_id`) REFERENCES `categories`(`id`)
// Jan 20 06:04:24 ip-172-31-41-198 cartique[8244]:  ┌───────────────────────────────────────────────────┐
