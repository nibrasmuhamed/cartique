package util

import "github.com/gofiber/fiber/v2"

func InternalServerErr(c *fiber.Ctx) error {
	return c.Status(500).JSON(fiber.Map{
		"message": "internal server error",
	})
}
