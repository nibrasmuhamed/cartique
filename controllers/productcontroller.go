package controllers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/google/uuid"
	"github.com/nibrasmuhamed/cartique/database"
	"github.com/nibrasmuhamed/cartique/models"
	"github.com/nibrasmuhamed/cartique/util"
	"gorm.io/gorm"
)

type ProductDB struct {
	DB *gorm.DB
}

func NewProductDB(p *gorm.DB) *ProductDB {
	return &ProductDB{p}
}

func (Pd *ProductDB) ShowAProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.OpenDataBase()
	defer database.CloseDatabase(db)
	r := db.QueryRow("SELECT id, name, price, category_id,quantity, specs from products where products.deleted_at is null AND id=?", id)
	a := models.ProductResp{}
	err := r.Scan(&a.ID, &a.Name, &a.Price, &a.Category_id, &a.Quantity, &a.Specs)
	if err != nil {
		fmt.Println("error while scanning ", err)
		return c.Status(400).JSON(fiber.Map{"message": "product id doesn't exist"})
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
	return c.Status(200).JSON(fiber.Map{"message": "success", "product": a})
}

func (Pd *ProductDB) ShowProducts(c *fiber.Ctx) error {
	db := database.OpenDataBase()
	defer database.CloseDatabase(db)
	p := []models.ProductRespHomeDemo{}
	// uc.DB.Model(models.Product{}).Preload("Images", "photo is not null").Select("images.photo").Find(&p)

	r, err := db.Query("SELECT products.id, name, price, category_id from products where products.deleted_at is null")
	if err != nil {
		fmt.Println("error is :", err)
		return c.Status(500).JSON(fiber.Map{"message": "internal server error"})
	}
	defer r.Close()
	for r.Next() {
		a := models.ProductRespHomeDemo{}
		err = r.Scan(&a.ID, &a.Name, &a.Price, &a.Category_id)
		if err != nil {
			fmt.Println("error while scanning ", err)
		}
		row := db.QueryRow("select photo from images where product_id = ? AND deleted_at is NULL LIMIT 1", a.ID)
		err := row.Scan(&a.Images)
		if err != nil {
			fmt.Println(err)
		}
		p = append(p, a)
	}
	return c.Status(200).JSON(fiber.Map{"message": "success", "product": p})
}

func (pd *AdminController) AddProduct(c *fiber.Ctx) error {
	p := models.Product{}
	x, _ := strconv.Atoi(c.FormValue("category_id"))
	p.Category_id = uint(x)
	p.Price, _ = strconv.Atoi(c.FormValue("price"))
	p.Quantity, _ = strconv.Atoi(c.FormValue("quantity"))
	p.Specs = c.FormValue("spec")
	p.Name = c.FormValue("name")
	form, _ := c.MultipartForm()
	if check := util.CheckFiles(p); !check {
		return c.Status(400).JSON(fiber.Map{"message": "required fields are empty"})
	}
	file := form.File["images"]
	if len(file) < 3 {
		return c.Status(400).JSON(fiber.Map{"message": "please upload atleast three images"})
	}
	for _, fileHeader := range file {
		x := uuid.New().String()
		c.SaveFile(fileHeader, fmt.Sprintf("./public/images/%s", x+".jpg"))
		u := models.Image{Product_id: p.ID, Photo: c.BaseURL() + "/images/" + x + ".jpg"}
		p.Images = append(p.Images, u)
		pd.DB.Create(&p.Images)
	}
	tx := pd.DB.Create(&p)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return c.Status(500).JSON(fiber.Map{"message": "cannot process the request now"})
	}
	return c.Status(200).JSON(fiber.Map{"message": "success"})
}

func (pd *AdminController) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var p models.Product
	pd.DB.Where("id = ?", id).Delete(&p)
	return c.Status(200).JSON(fiber.Map{"message": "success"})
}

func (pd *AdminController) EditProduct(c *fiber.Ctx) error {
	p, img := models.Product{}, models.Image{Product_id: 0}
	pd.DB.First(&p, "id = ?", c.Params("id"))
	tx := pd.DB.Where("product_id = ?", c.Params("id")).Unscoped().Delete(&img)
	if tx.Error != nil {
		fmt.Println("txerr", tx.Error)
	}
	x, _ := strconv.Atoi(c.FormValue("category_id"))
	p.Category_id = uint(x)
	p.Price, _ = strconv.Atoi(c.FormValue("price"))
	p.Quantity, _ = strconv.Atoi(c.FormValue("quantity"))
	p.Specs = c.FormValue("spec")
	p.Name = c.FormValue("name")
	form, _ := c.MultipartForm()

	if check := util.CheckFiles(p); !check {
		return c.Status(400).JSON(fiber.Map{"message": "required fields are empty"})
	}
	file := form.File["images"]
	if len(file) < 3 {
		return c.Status(400).JSON(fiber.Map{"message": "please upload atleast three images"})
	}

	for _, fileHeader := range file {
		x := uuid.New().String()
		c.SaveFile(fileHeader, fmt.Sprintf("./public/images/%s", x+".jpg"))
		u := models.Image{Product_id: p.ID, Photo: c.BaseURL() + "/images/" + x + ".jpg"}
		p.Images = append(p.Images, u)
		pd.DB.Save(&p.Images)
	}
	tx = pd.DB.Save(&p)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return c.Status(500).JSON(fiber.Map{"message": "cannot process the request now"})
	}
	return c.Status(200).JSON(fiber.Map{"message": "success"})
}

func (uc *AdminController) Logout(c *fiber.Ctx) error {
	t, _, err := util.GetAuthRef(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "user token not found"})
	}
	id := util.GetidfromToken(t)
	var u models.Admin
	uc.DB.Where("id = ?", id).First(&u)
	u.Refresh_token = ""
	uc.DB.Save(&u)
	t, err = util.GenerateJWT(int(id), "admin", 0)
	if err != nil {
		fmt.Println("something")
	}
	return c.Status(200).JSON(fiber.Map{"message": "logged out",
		"access_token": t})
}
