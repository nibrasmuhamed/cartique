package controllers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/google/UUID"
	"github.com/nibrasmuhamed/cartique/database"
	"github.com/nibrasmuhamed/cartique/models"
	"github.com/nibrasmuhamed/cartique/util"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) *UserController {
	return &UserController{DB}
}

func (uc *UserController) RegisterUser(c *fiber.Ctx) error {
	u := new(models.UserRegister)
	if err := c.BodyParser(u); err != nil {
		log.Println("bad time on parsing json", err)
		c.SendStatus(500)
		return c.JSON(fiber.Map{"message": "internal server error"})
	}
	if validEmail := util.ValidMailAddress(u.Email); !validEmail {
		return c.Status(400).JSON(fiber.Map{"message": "enter a valid email"})
	}
	if u.Password != u.Password1 {
		return c.Status(400).JSON(fiber.Map{"message": "password do not match"})
	}
	var d int64
	uc.DB.Model(models.User{}).Where("email = ?", u.Email).Count(&d)
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
	uc.DB.Create(&user)
	c.SendStatus(200)
	return c.JSON(fiber.Map{"message": "user created succesfully"})
}

func (uc *UserController) LoginUser(c *fiber.Ctx) error {
	var u models.UserLogin
	if err := c.BodyParser(&u); err != nil {
		c.SendStatus(500)
		return c.JSON(fiber.Map{"message": "internal server error"})
	}
	// db := database.OpenDb()
	// defer database.CloseDb(db)
	user := models.User{}
	err := uc.DB.Where("email = ?", u.Email).First(&user).Error
	if err != nil {
		fmt.Println("error is ", err)
		return c.Status(500).JSON(fiber.Map{"message": "internal server error"})
	}
	if user.Email == "" {
		c.SendStatus(400)
		return c.JSON(fiber.Map{"message": "email or password is incorrect"})
	}
	if verified := util.VerifyPassword(u.Password, user.Password); !verified {
		c.SendStatus(401)
		return c.JSON(fiber.Map{"message": "email or password is incorrect"})
	}
	if user.Blocked {
		return c.Status(400).JSON(fiber.Map{"message": "your account is suspended by cartique"})
	}
	token, err := util.GenerateJWT(int(user.ID), "user", 100000)
	if err != nil {
		log.Println(err)
	}
	uuidv4, _ := uuid.NewRandom()
	uc.DB.Model(&user).Update("refresh_token", uuidv4)
	res := models.UserResponse{Id: int(user.ID), Username: user.Username, Email: user.Email}
	return c.Status(200).JSON(fiber.Map{"details": res, "message": "success", "access_token": token, "refresh_token": uuidv4})
}

func (uc *UserController) VerifyUser(c *fiber.Ctx) error {
	t := c.Get("Authorization")
	t = t[7:]
	id := util.GetidfromToken(t)
	var p models.User
	uc.DB.Where("id = ?", id).First(&p)
	if x := util.SendOtp(p.Phone); !x {
		return c.Status(500).JSON(fiber.Map{"message": "internal sever error"})
	}
	return c.JSON(fiber.Map{"message": "success"})
}

func (uc *UserController) VerifyUserOtp(c *fiber.Ctx) error {
	otp := c.Params("id")
	t := c.Get("Authorization")
	t = t[7:]
	id := util.GetidfromToken(t)
	var p models.User
	uc.DB.Where("id = ?", id).First(&p)
	verified := util.CheckOtp(p.Phone, otp)
	if !verified {
		return c.Status(400).JSON(fiber.Map{"message": "cannot verify"})
	}
	p.Verified = true
	uc.DB.Save(&p)
	return c.Status(200).JSON(fiber.Map{"message": "success"})
}

func (uc *UserController) ShowProducts(c *fiber.Ctx) error {
	db := database.OpenDataBase()
	defer database.CloseDatabase(db)
	p := []models.ProductRespHome{}
	// uc.DB.Model(models.Product{}).Preload("Images", "photo is not null").Select("images.photo").Find(&p)

	r, err := db.Query("SELECT id, name, price, category_id from products where products.deleted_at is null")
	if err != nil {
		fmt.Println("error is :", err)
		return c.Status(500).JSON(fiber.Map{"message": "internal server error"})
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
	return c.Status(200).JSON(fiber.Map{"message": "success", "product": p})
}

func (uc *UserController) RefreshToken(c *fiber.Ctx) error {
	t := c.Get("Authorization")
	r := c.Get("refresh_token")
	if t == "" || r == "" {
		return c.Status(404).JSON(fiber.Map{"message": "access token or refresh token not found"})
	}
	t = t[7:]
	var u models.User
	verified, _ := util.CheckToken(t)
	if !verified {
		u.Refresh_token = ""
		uc.DB.Save(&u)
		return c.Status(401).JSON(fiber.Map{"message": "user not authorized"})
	}
	id := int(util.GetidfromToken(string(t)))
	uc.DB.Where("id = ?", id).First(&u)
	if r != u.Refresh_token {
		u.Refresh_token = ""
		uc.DB.Save(&u)
		return c.Status(401).JSON(fiber.Map{"message": "your refresh token has been edited"})
	}
	d := util.CompareTime(u.UpdatedAt)
	if d > 5 {
		u.Refresh_token = ""
		uc.DB.Save(&u)
		return c.Status(401).JSON(fiber.Map{"message": "your refresh token has been expired"})
	}
	token, _ := util.GenerateJWT(int(id), "user", 5)
	uuidv4, _ := uuid.NewRandom()
	uc.DB.Model(&u).Update("refresh_token", uuidv4)
	return c.Status(200).JSON(fiber.Map{"message": "success",
		"access_token": token, "refresh_token": uuidv4})
}

func (uc *UserController) Logout(c *fiber.Ctx) error {
	t, _, err := util.GetAuthRef(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "user token not found"})
	}
	id := util.GetidfromToken(t)
	var u models.User
	uc.DB.Where("id = ?", id).First(&u)
	u.Refresh_token = ""
	uc.DB.Save(&u)
	t, err = util.GenerateJWT(int(id), "user", 0)
	if err != nil {
		fmt.Println("something")
	}
	return c.Status(200).JSON(fiber.Map{"message": "logged out",
		"access_token": t})
}

func (uc *UserController) AddToCart(c *fiber.Ctx) error {
	productID, _ := strconv.Atoi(c.Params("id"))
	t := c.Get("Authorization")
	if t == "" {
		return c.Status(404).JSON(fiber.Map{"message": "access token or refresh token not found"})
	}
	t = t[7:]
	userID := util.GetidfromToken(t)
	var cart models.Cart
	tx := uc.DB.Where("product_id=? AND user_id=?", productID, userID).First(&cart)
	if tx.Error != gorm.ErrRecordNotFound {
		cart.Quantity += 1
		uc.DB.Save(&cart)
		return c.Status(200).JSON(fiber.Map{"message": "success"})
	}
	cart = models.Cart{ProductID: uint(productID), UserID: uint(userID), Quantity: 1}
	uc.DB.Create(&cart)
	return c.Status(200).JSON(fiber.Map{"message": "success"})
}

func (uc *UserController) ShowCart(c *fiber.Ctx) error {
	t := c.Get("Authorization")
	t = t[7:]
	userID := util.GetidfromToken(t)
	db := database.OpenDataBase()
	defer database.CloseDatabase(db)
	prods := []int{}
	rows, err := db.Query("SELECT product_id FROM carts WHERE user_id=?", userID)
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	if err != nil {
		fmt.Println("error is :", err)
	}
	for rows.Next() {
		var x int
		rows.Scan(&x)
		prods = append(prods, x)
	}
	x := util.FindProducts(prods)
	return c.Status(200).JSON(fiber.Map{"message": "success",
		"products": x})

}

func (uc *UserController) AddAddress(c *fiber.Ctx) error {
	userId := int(c.Locals("userid").(float64))
	ad := models.Address{}
	if err := c.BodyParser(&ad); err != nil {
		fmt.Println(err)
		return util.InternalServerErr(c)
	}
	ad.UserID = uint(userId)
	uc.DB.Create(&ad)
	return c.Status(200).JSON(fiber.Map{"message": "success"})
}

func (uc *UserController) ShowAddress(c *fiber.Ctx) error {
	userId := int(c.Locals("userid").(float64))
	db := database.OpenDataBase()
	defer database.CloseDatabase(db)
	ads := []models.AddressResp{}
	rows, err := db.Query("SELECT id, user_id, home,city, disctrict,state, pin FROM addresses WHERE user_id=?", userId)
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var ad models.AddressResp
		if err := rows.Scan(&ad.ID, &ad.UserID, &ad.Home, &ad.City, &ad.District, &ad.State, &ad.Pin); err != nil {
			fmt.Println(err)
		}
		ads = append(ads, ad)
	}
	return c.Status(200).JSON(fiber.Map{"message": "success",
		"addresses": ads,
	})
}

func (uc *UserController) EditUser(c *fiber.Ctx) error {
	userId := int(c.Locals("userid").(float64))
	var u, udb models.User
	if err := c.BodyParser(&u); err != nil {
		fmt.Println(err)
	}
	uc.DB.First(&udb, userId)
	udb.Email = u.Email
	udb.Phone = u.Phone
	udb.Username = u.Username
	fmt.Println(udb)
	uc.DB.Save(&udb)
	return c.Status(200).JSON(fiber.Map{"message": "success"})
}

func (uc *UserController) Showwishlist(c *fiber.Ctx) error {
	t := c.Get("Authorization")
	t = t[7:]
	userID := util.GetidfromToken(t)
	db := database.OpenDataBase()
	defer database.CloseDatabase(db)
	prods := []int{}
	rows, err := db.Query("SELECT product_id FROM wishlists WHERE user_id=?", userID)
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	if err != nil {
		fmt.Println("error is :", err)
	}
	for rows.Next() {
		var x int
		rows.Scan(&x)
		prods = append(prods, x)
	}
	x := util.FindProducts(prods)
	return c.Status(200).JSON(fiber.Map{"message": "success",
		"products": x})

}

func (uc *UserController) AddToWishlist(c *fiber.Ctx) error {
	productID, _ := strconv.Atoi(c.Params("id"))
	t := c.Get("Authorization")
	if t == "" {
		return c.Status(404).JSON(fiber.Map{"message": "access token or refresh token not found"})
	}
	t = t[7:]
	userID := util.GetidfromToken(t)
	var cart models.Wishlist
	tx := uc.DB.Where("product_id=? AND user_id=?", productID, userID).First(&cart)
	if tx.Error != gorm.ErrRecordNotFound {
		cart.Quantity += 1
		uc.DB.Save(&cart)
		return c.Status(200).JSON(fiber.Map{"message": "success"})
	}
	cart = models.Wishlist{ProductID: uint(productID), UserID: uint(userID), Quantity: 1}
	uc.DB.Create(&cart)
	return c.Status(200).JSON(fiber.Map{"message": "success"})
}

func (uc *UserController) RemoveFromWishList(c *fiber.Ctx) error {
	productID, _ := strconv.Atoi(c.Params("id"))
	t := c.Get("Authorization")
	if t == "" {
		return c.Status(404).JSON(fiber.Map{"message": "access token or refresh token not found"})
	}
	t = t[7:]
	userID := util.GetidfromToken(t)
	var cart models.Wishlist
	uc.DB.Where("product_id=? AND user_id=?", productID, userID).Delete(&cart)
	return c.Status(200).JSON(fiber.Map{"message": "success"})
}

func (uc *UserController) RemoveFromCart(c *fiber.Ctx) error {
	productID, _ := strconv.Atoi(c.Params("id"))
	t := c.Get("Authorization")
	if t == "" {
		return c.Status(404).JSON(fiber.Map{"message": "access token or refresh token not found"})
	}
	t = t[7:]
	userID := util.GetidfromToken(t)
	var cart models.Cart
	uc.DB.Where("product_id=? AND user_id=?", productID, userID).Delete(&cart)
	return c.Status(200).JSON(fiber.Map{"message": "success"})
}
