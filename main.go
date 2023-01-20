package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/nibrasmuhamed/cartique/controllers"
	"github.com/nibrasmuhamed/cartique/database"
	"github.com/nibrasmuhamed/cartique/routes"
)

var (
	UC *controllers.UserController
	AC *controllers.AdminController
	ar *routes.AdminRouter
	ur *routes.UserRouter
	Pc *controllers.ProductDB
	Pr *routes.ProductRouter
)

func init() {
	x := database.InitDB()
	AC = controllers.NewAdminController(x)
	UC = controllers.NewUserController(x)
	Pc = controllers.NewProductDB(x)
	ur = routes.NewUserRouter(UC)
	ar = routes.NewAdminRouter(*AC)
	Pr = routes.NewProductRoter(Pc)
}

func main() {

	app := fiber.New()
	// app.Use(csrf.New())
	app.Static("/", "./public")
	// app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	products := app.Group("/products")
	Pr.Routes(products)
	admin := app.Group("/admin")
	ar.AdminRoute(admin)
	user := app.Group("/user")
	ur.Routes(user)

	log.Fatal(app.Listen("0.0.0.0:60000"))
}
