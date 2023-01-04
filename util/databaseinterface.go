package util

import "github.com/gofiber/fiber/v2"

type DB interface {
	RegisterAdmin(c *fiber.Ctx) error
}
