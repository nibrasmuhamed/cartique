package controllers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/google/UUID"
	"github.com/nibrasmuhamed/cartique/database"
	"github.com/nibrasmuhamed/cartique/models"
	"github.com/nibrasmuhamed/cartique/util"
)

func AddProduct(c *fiber.Ctx) error {
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
	db := database.OpenDb()
	defer database.CloseDb(db)
	for _, fileHeader := range file {
		x := uuid.New().String()
		c.SaveFile(fileHeader, fmt.Sprintf("./public/images/%s", x+".jpg"))
		u := models.Image{Product_id: p.ID, Photo: c.BaseURL() + "/images/" + x + ".jpg"}
		p.Images = append(p.Images, u)
		db.Create(&p.Images)
	}
	tx := db.Create(&p)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return c.Status(500).JSON(fiber.Map{"message": "cannot process the request now"})
	}
	return c.Status(200).JSON(fiber.Map{"message": "success"})
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	db := database.OpenDb()
	defer database.CloseDb(db)
	var p models.Product
	db.Where("id = ?", id).Delete(&p)
	return c.Status(200).JSON(fiber.Map{"message": "success"})
}

func EditProduct(c *fiber.Ctx) error {
	p, img := models.Product{}, models.Image{Product_id: 0}
	db := database.OpenDb()
	defer database.CloseDb(db)
	db.First(&p, "id = ?", c.Params("id"))
	tx := db.Where("product_id = ?", c.Params("id")).Unscoped().Delete(&img)
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
		db.Save(&p.Images)
	}
	tx = db.Save(&p)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return c.Status(500).JSON(fiber.Map{"message": "cannot process the request now"})
	}
	return c.Status(200).JSON(fiber.Map{"message": "success"})
}
