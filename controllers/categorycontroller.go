package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nibrasmuhamed/cartique/database"
	"github.com/nibrasmuhamed/cartique/models"
)

func AddCategory(c *fiber.Ctx) error {
	ct := &models.Category{}
	err := c.BodyParser(&ct)
	if err != nil {
		log.Println(err)
		return c.SendStatus(500)

	}
	db := database.OpenDb()
	defer database.CloseDb(db)
	tx := db.Create(&ct)
	if tx.Error != nil {
		return c.Status(400).JSON(fiber.Map{"message": "catergory exist"})
	}
	return c.Status(200).JSON(fiber.Map{"message": "success"})
}
