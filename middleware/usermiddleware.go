package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nibrasmuhamed/cartique/util"
)

func UserMiddleware(c *fiber.Ctx) error {
	t := c.Get("Authorization")
	if t == "" {
		return c.Status(404).JSON(fiber.Map{"message": "token not found"})
	}
	// bearer
	t = t[7:]
	verified, _ := util.CheckToken(t)
	if !verified {
		return c.Status(401).JSON(fiber.Map{"message": "user not authorized"})
	}
	userID := util.GetidfromToken(t)
	// return c.Status(200).JSON(fiber.Map{"message": "success"})
	c.Locals("userid", userID)
	return c.Next()
}
