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
