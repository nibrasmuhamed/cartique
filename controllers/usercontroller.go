package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nibrasmuhamed/cartique/database"
	"github.com/nibrasmuhamed/cartique/models"
	"github.com/nibrasmuhamed/cartique/util"
)

func RegisterUser(c *fiber.Ctx) error {
	u := new(models.UserRegister)
	if err := c.BodyParser(u); err != nil {
		log.Println("bad time on parsing json", err)
		c.SendStatus(500)
		return c.JSON(fiber.Map{"message": "internal server error"})
	}
	db := database.OpenDb()
	defer database.CloseDb(db)
	var d int64
	db.Model(models.User{}).Where("email = ?", u.Email).Count(&d)
	if d >= 1 {
		c.SendStatus(400)
		return c.JSON(fiber.Map{"message": "user with this email already exist"})
	}
	if valid := util.Password(u.Password); !valid {
		c.SendStatus(400)
		return c.JSON(fiber.Map{"message": "password is not strong"})
	}

	hash := util.HashPassword(u.Password)
	user := models.User{Username: u.Name, Email: u.Email, Phone: u.Phone, Password: hash, Verified: false}
	db.Create(&user)
	c.SendStatus(200)
	return c.JSON(fiber.Map{"message": "user created succesfully"})
}
