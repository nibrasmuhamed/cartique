package util

import "github.com/gofiber/fiber/v2"

func InternalServerErr(c *fiber.Ctx) error {
	return c.Status(500).JSON(fiber.Map{
		"message": "internal server error",
	})
}

func BadReq(c *fiber.Ctx, message string) error {
	return c.Status(500).JSON(fiber.Map{
		"message": message,
	})
}
