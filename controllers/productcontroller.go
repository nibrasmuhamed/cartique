package controllers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/google/UUID"
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

func (Pd *ProductDB) ShowProducts(c *fiber.Ctx) error {
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

func (pd *AdminController) AddProduct(c *fiber.Ctx) error {
	p := models.Product{}
	p.Category_id, _ = strconv.Atoi(c.FormValue("category_id"))
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
	p.Category_id, _ = strconv.Atoi(c.FormValue("category_id"))
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
