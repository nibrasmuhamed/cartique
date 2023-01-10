package controllers

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/google/UUID"
	"github.com/nibrasmuhamed/cartique/models"
	"github.com/nibrasmuhamed/cartique/util"
	"gorm.io/gorm"
)

type AdminController struct {
	DB *gorm.DB
}

func NewAdminController(DB *gorm.DB) *AdminController {
	return &AdminController{DB}
}

func (ac *AdminController) EditUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var u models.User
	ac.DB.Where("id = ?", id).First(&u)
	if u.Email == "" {
		return c.Status(404).JSON(fiber.Map{"message": "cannot find the user"})
	}
	u.Blocked = true
	u.Refresh_token = ""
	ac.DB.Save(&u)
	return c.Status(200).JSON(fiber.Map{"message": "user blocked succesfully"})
}

func (ac *AdminController) UnBlockUsers(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var u models.User
	ac.DB.Where("id = ?", id).First(&u)
	if u.Email == "" {
		return c.Status(404).JSON(fiber.Map{"message": "cannot find the user"})
	}
	u.Blocked = false
	u.Refresh_token = ""
	ac.DB.Save(&u)
	return c.Status(200).JSON(fiber.Map{"message": "user unblocked succesfully"})
}

func (ac *AdminController) BlockUsers(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var u models.User
	ac.DB.Where("id = ?", id).First(&u)
	if u.Email == "" {
		return c.Status(404).JSON(fiber.Map{"message": "cannot find the user"})
	}
	u.Blocked = true
	u.Refresh_token = ""
	ac.DB.Save(&u)
	return c.Status(200).JSON(fiber.Map{"message": "user blocked succesfully"})
}

func (ac *AdminController) LoginAdmin(c *fiber.Ctx) error {
	a := new(models.AdminLogin)
	var admindb models.Admin
	if err := c.BodyParser(a); err != nil {
		log.Println("error while parsing data: ", err)
		return err
	}
	tx := ac.DB.Where("email = ?", a.Email).First(&admindb)
	if tx.Error == gorm.ErrRecordNotFound {
		c.SendStatus(400)
		return c.JSON(fiber.Map{"message": "no user found associated with this email"})
	}
	ok := util.VerifyPassword(a.Password, admindb.Password)
	if !ok {
		c.SendStatus(403)
		return c.JSON(fiber.Map{"message": "username or password is incorrect"})
	}
	token, _ := util.GenerateJWT(admindb.Id, "admin", 5)
	ar := models.AdminResponse{Id: admindb.Id, Email: admindb.Email, Username: admindb.Username}
	c.SendStatus(200)
	return c.JSON(fiber.Map{"message": "login successful",
		"details":      ar,
		"access_token": token})
}

func (ac *AdminController) RefressTokenAdmin(c *fiber.Ctx) error {

	t := c.Get("Authorization")
	r := c.Get("refresh_token")
	if t == "" || r == "" {
		return c.Status(404).JSON(fiber.Map{"message": "access token or refresh token not found"})
	}
	t = t[7:]
	var a models.Admin
	verified, _ := util.CheckToken(t)
	if !verified {
		a.Refresh_token = ""
		ac.DB.Save(&a)
		return c.Status(401).JSON(fiber.Map{"message": "user not authorized"})
	}
	id := int(util.GetidfromToken(string(t)))
	ac.DB.Where("id = ?", id).First(&a)
	if r != a.Refresh_token {
		a.Refresh_token = ""
		ac.DB.Save(&a)
		return c.Status(401).JSON(fiber.Map{"message": "your refresh token has been edited"})
	}
	token, _ := util.GenerateJWT(int(id), "admin", 5)
	uuidv4, _ := uuid.NewRandom()
	ac.DB.Model(&a).Update("refresh_token", uuidv4)
	return c.Status(200).JSON(fiber.Map{"message": "success",
		"access_token": token, "refresh_token": uuidv4})
}

func (ac *AdminController) RegisterAdmin(c *fiber.Ctx) error {
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
	// checking whether an admin exist with this email or not
	tx := ac.DB.Where("email = ?", a.Email).First(&adb)
	if tx.Error != gorm.ErrRecordNotFound {
		c.SendStatus(400)
		return c.JSON(fiber.Map{"message": "admin with this email already exist"})
	}
	u := models.Admin{Username: a.Username, Email: a.Email, Password: util.HashPassword(a.Password)}
	ac.DB.Create(&u)
	c.SendStatus(200)
	return c.JSON(fiber.Map{"message": "success"})
}

func (ac *AdminController) ViewUsers(c *fiber.Ctx) error {
	l, _ := strconv.Atoi(c.Query("limit"))
	o, _ := strconv.Atoi(c.Query("offset"))

	var users []models.UserResponse
	ac.DB.Model(models.User{}).Limit(l).Offset(o).Find(&users)
	return c.Status(200).JSON(fiber.Map{"message": "success", "users": users})
}
