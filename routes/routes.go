package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/nibrasmuhamed/cartique/controllers"
	"github.com/nibrasmuhamed/cartique/middleware"
)

func Routes() *fiber.App {
	app := fiber.New()
	admin := app.Group("/admin")
	admin.Post("/register", controllers.RegisterAdmin)
	admin.Post("/login", controllers.LoginAdmin)
	admin.Get("/user_managment", middleware.AdminMiddleware, controllers.ViewUsers)
	admin.Get("/user_managment/block/:id", middleware.AdminMiddleware, controllers.BlockUsers)
	admin.Get("/user_managment/unblock/:id", middleware.AdminMiddleware, controllers.UnBlockUsers)
	admin.Get("/refresh", controllers.Refresh_token_admin)
	user := app.Group("/user")
	user.Post("/register", controllers.RegisterUser)
	user.Post("/login", controllers.LoginUser)
	user.Get("/verify", middleware.UserMiddleware, controllers.VerifyUser)
	user.Get("/verify/:id", middleware.UserMiddleware, controllers.VerifyUserOtp)
	user.Get("/refresh", controllers.RefreshToken)
	user.Get("/logout", controllers.Logout)
	app.Static("/", "./public")
	app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	return app
}
