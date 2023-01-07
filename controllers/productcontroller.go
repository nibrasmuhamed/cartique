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
	file := form.File["images"]
	if check := util.CheckFiles(p); !check {
		return c.Status(400).JSON(fiber.Map{"message": "required fields are empty"})
	}
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
	}
	db.Create(&p)
	return c.Status(200).JSON(fiber.Map{"message": "success"})
}
