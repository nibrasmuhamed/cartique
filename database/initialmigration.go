package database

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/nibrasmuhamed/cartique/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type db *gorm.DB

func InitDB() {
	db, err := gorm.Open(mysql.Open(os.Getenv("DB")), &gorm.Config{})
	if err != nil {
		log.Println("error in connecting database : ", err)
	}
	db.AutoMigrate(models.User{})
	db.AutoMigrate(models.Admin{})

	defer CloseDb(db)
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
