package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	utils "github.com/gofiber/utils"
)

type UserControllerFake struct {
	Calls map[string][][]interface{}
}

func TestShowProducts(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendStatus(400)
	})

	resp, err := app.Test(httptest.NewRequest("GET", "/test", nil))

	utils.AssertEqual(t, nil, err, "app.Test")
	utils.AssertEqual(t, 200, resp.StatusCode, "Status code")
}
