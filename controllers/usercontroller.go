package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/google/UUID"
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

func LoginUser(c *fiber.Ctx) error {
	var u models.UserLogin
	if err := c.BodyParser(&u); err != nil {
		c.SendStatus(500)
		return c.JSON(fiber.Map{"message": "internal server error"})
	}
	db := database.OpenDb()
	defer database.CloseDb(db)
	user := models.User{}
	db.Where("email = ?", u.Email).First(&user)
	if user.Email == "" {
		c.SendStatus(400)
		return c.JSON(fiber.Map{"message": "email or password is incorrect"})
	}
	if verified := util.VerifyPassword(u.Password, user.Password); !verified {
		c.SendStatus(401)
		return c.JSON(fiber.Map{"message": "email or password is incorrect"})
	}
	token, err := util.GenerateJWT(int(user.ID), "user")
	if err != nil {
		log.Println(err)
	}
	uuidv4, _ := uuid.NewRandom()
	db.Model(&user).Update("refresh_token", uuidv4)
	res := models.UserResponse{Id: int(user.ID), Username: user.Username, Email: user.Email}
	return c.Status(200).JSON(fiber.Map{"details": res, "message": "success", "access_token": token, "refresh_token": uuidv4})
}

func VerifyUser(c *fiber.Ctx) error {
	t := c.Get("Authorization")
	t = t[7:]
	id := util.GetidfromToken(t)
	db := database.OpenDb()
	defer database.CloseDb(db)
	var p models.User
	db.Where("id = ?", id).First(&p)
	if x := util.SendOtp(p.Phone); !x {
		return c.Status(500).JSON(fiber.Map{"message": "internal sever error"})
	}
	return c.JSON(fiber.Map{"message": "success"})
}

func VerifyUserOtp(c *fiber.Ctx) error {
	otp := c.Params("id")
	t := c.Get("Authorization")
	t = t[7:]
	id := util.GetidfromToken(t)
	db := database.OpenDb()
	defer database.CloseDb(db)
	var p models.User
	db.Where("id = ?", id).First(&p)
	verified := util.CheckOtp(p.Phone, otp)
	if !verified {
		return c.Status(400).JSON(fiber.Map{"message": "cannot verify"})
	}
	p.Verified = true
	db.Save(&p)
	return c.Status(200).JSON(fiber.Map{"message": "success"})
}
