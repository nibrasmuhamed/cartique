package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nibrasmuhamed/cartique/util"
)

func AdminMiddleware(c *fiber.Ctx) error {
	t := c.Get("Authorization")
	if t == "" {
		return c.Status(404).JSON(fiber.Map{"message": "token not found"})
	}
	t = t[7:]
	verified, _ := util.CheckTokenAdmin(t)
	if !verified {
		return c.Status(401).JSON(fiber.Map{"message": "user not authorized"})
	}
	// return c.Status(200).JSON(fiber.Map{"message": "success"})
	return c.Next()
}
