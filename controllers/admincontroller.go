package controllers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nibrasmuhamed/cartique/database"
	"github.com/nibrasmuhamed/cartique/models"
	"github.com/nibrasmuhamed/cartique/util"
	"gorm.io/gorm"
)

func LoginAdmin(c *fiber.Ctx) error {
	a := new(models.AdminLogin)
	var admindb models.Admin
	if err := c.BodyParser(a); err != nil {
		log.Println("error while parsing data: ", err)
		return err
	}
	db := database.OpenDb()
	defer database.CloseDb(db)
	tx := db.Where("email = ?", a.Email).First(&admindb)
	if tx.Error == gorm.ErrRecordNotFound {
		c.SendStatus(400)
		return c.JSON(fiber.Map{"message": "no user found associated with this email"})
	}
	ok := util.VerifyPassword(a.Password, admindb.Password)
	if !ok {
		c.SendStatus(403)
		return c.JSON(fiber.Map{"message": "username or password is incorrect"})
	}
	c.SendStatus(200)
	return c.JSON(fiber.Map{"message": "login successful"})
}

func RegisterAdmin(c *fiber.Ctx) error {
	a := new(models.AdminRegister)
	adb := new(models.Admin)
	// parsing body in to the admin register struct.
	if err := c.BodyParser(a); err != nil {
		log.Println("error while parsing body : ", err)
	}
	// checking strong password.
	if valid := util.Password(a.Password); !valid {
		c.SendStatus(400)
		return c.JSON(fiber.Map{"message": "password is not strong"})
	}
	// database opening and defering close
	db := database.OpenDb()
	defer database.CloseDb(db)
	// checking whether an admin exist with this email or not
	tx := db.Where("email = ?", a.Email).First(&adb)
	if tx.Error != gorm.ErrRecordNotFound {
		c.SendStatus(400)
		return c.JSON(fiber.Map{"message": "admin with this email already exist"})
	}
	u := models.Admin{Username: a.Username, Email: a.Email, Password: util.HashPassword(a.Password)}
	db.Create(&u)
	c.SendStatus(200)
	return c.JSON(fiber.Map{"message": "success"})
}

func ViewUsers(c *fiber.Ctx) error {
	l, _ := strconv.Atoi(c.Query("limit"))
	o, _ := strconv.Atoi(c.Query("offset"))

	var users []models.UserResponse
	db := database.OpenDb()
	defer database.CloseDb(db)
	db.Model(models.User{}).Limit(l).Offset(o).Find(&users)
	return c.Status(200).JSON(fiber.Map{"message": "success", "users": users})
}
