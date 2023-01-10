package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nibrasmuhamed/cartique/controllers"
	"github.com/nibrasmuhamed/cartique/database"
	"github.com/nibrasmuhamed/cartique/routes"
)

var (
	UC *controllers.UserController
	a  *routes.UserRouter
)

func init() {
	x := database.InitDB()
	UC = controllers.NewUserController(x)
	a = routes.NewUserRouter(UC)
}

func main() {

	app := fiber.New()
	app.Static("/images", "./public/images")
	// app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))
	// app.Use(logger.New(logger.Config{
	// 	Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	// }))
	admin := app.Group("/admin")
	a.AdminRoute(admin)
	user := app.Group("/user")
	a.Routes(user)

	log.Fatal(app.Listen(":8000"))
}
