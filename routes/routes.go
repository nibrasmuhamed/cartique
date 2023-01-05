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
	admin.Get("/users", controllers.ViewUsers)

	user := app.Group("/user")
	user.Post("/register", controllers.RegisterUser)
	user.Post("/login", controllers.LoginUser)
	user.Get("/verify", middleware.UserMiddleware, controllers.VerifyUser)
	user.Get("/verify/:id", middleware.UserMiddleware, controllers.VerifyUserOtp)

	app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))

	return app
}
