package controllers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/google/UUID"
	"github.com/nibrasmuhamed/cartique/database"
	"github.com/nibrasmuhamed/cartique/models"
	"github.com/nibrasmuhamed/cartique/util"
	"gorm.io/gorm"
)

func UnBlockUsers(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	db := database.OpenDb()
	defer database.CloseDb(db)
	var u models.User
	db.Where("id = ?", id).First(&u)
	if u.Email == "" {
		return c.Status(404).JSON(fiber.Map{"message": "cannot find the user"})
	}
	u.Blocked = false
	u.Refresh_token = ""
	db.Save(&u)
	return c.Status(200).JSON(fiber.Map{"message": "user unblocked succesfully"})
}

func BlockUsers(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	db := database.OpenDb()
	defer database.CloseDb(db)
	var u models.User
	db.Where("id = ?", id).First(&u)
	if u.Email == "" {
		return c.Status(404).JSON(fiber.Map{"message": "cannot find the user"})
	}
	u.Blocked = true
	u.Refresh_token = ""
	db.Save(&u)
	return c.Status(200).JSON(fiber.Map{"message": "user blocked succesfully"})
}

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
	token, _ := util.GenerateJWT(admindb.Id, "admin")
	ar := models.AdminResponse{Id: admindb.Id, Email: admindb.Email, Username: admindb.Username}
	c.SendStatus(200)
	return c.JSON(fiber.Map{"message": "login successful",
		"details":      ar,
		"access_token": token})
}

func Refresh_token_admin(c *fiber.Ctx) error {

	t := c.Get("Authorization")
	r := c.Get("refresh_token")
	if t == "" || r == "" {
		return c.Status(404).JSON(fiber.Map{"message": "access token or refresh token not found"})
	}
	t = t[7:]
	var a models.Admin
	db := database.OpenDb()
	defer database.CloseDb(db)
	verified, _ := util.CheckToken(t)
	if !verified {
		a.Refresh_token = ""
		db.Save(&a)
		return c.Status(401).JSON(fiber.Map{"message": "user not authorized"})
	}
	id := int(util.GetidfromToken(string(t)))
	db.Where("id = ?", id).First(&a)
	if r != a.Refresh_token {
		a.Refresh_token = ""
		db.Save(&a)
		return c.Status(401).JSON(fiber.Map{"message": "your refresh token has been edited"})
	}
	token, _ := util.GenerateJWT(int(id), "user")
	uuidv4, _ := uuid.NewRandom()
	db.Model(&a).Update("refresh_token", uuidv4)
	return c.Status(200).JSON(fiber.Map{"message": "success",
		"access_token": token, "refresh_token": uuidv4})
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
