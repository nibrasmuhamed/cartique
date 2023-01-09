package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/nibrasmuhamed/cartique/controllers"
	"github.com/nibrasmuhamed/cartique/middleware"
)

func Routes() *fiber.App {
	app := fiber.New()
	admin := app.Group("/admin")
	admin.Post("/register", controllers.RegisterAdmin)
	admin.Post("/login", controllers.LoginAdmin)
	admin.Get("/refresh", controllers.Refresh_token_admin)

	userManagment := admin.Group("user_managment")
	userManagment.Get("/", middleware.AdminMiddleware, controllers.ViewUsers)
	userManagment.Get("/block/:id", middleware.AdminMiddleware, controllers.BlockUsers)
	userManagment.Get("/unblock/:id", middleware.AdminMiddleware, controllers.UnBlockUsers)
	product := admin.Group("/products")
	product.Post("/add_category", controllers.AddCategory)
	product.Post("/add_product", controllers.AddProduct)
	product.Post("/delete_product/:id", controllers.DeleteProduct)

	user := app.Group("/user")
	user.Post("/register", controllers.RegisterUser)
	user.Post("/login", controllers.LoginUser)
	user.Get("/verify", middleware.UserMiddleware, controllers.VerifyUser)
	user.Get("/verify/:id", middleware.UserMiddleware, controllers.VerifyUserOtp)
	user.Get("/refresh", controllers.RefreshToken)
	user.Get("/logout", controllers.Logout)
	user.Get("/", controllers.ShowProducts)

	app.Static("/images", "./public/images")
	// app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	return app
}
