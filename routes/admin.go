package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nibrasmuhamed/cartique/controllers"
	"github.com/nibrasmuhamed/cartique/middleware"
)

func (u *UserRouter) AdminRoute(admin fiber.Router) {
	admin.Post("/register", controllers.RegisterAdmin)
	admin.Post("/login", u.UserRoute.LoginAdmin)
	admin.Get("/refresh", controllers.Refresh_token_admin)

	userManagment := admin.Group("user_managment")
	userManagment.Get("/", middleware.AdminMiddleware, controllers.ViewUsers)
	userManagment.Get("/block/:id", middleware.AdminMiddleware, controllers.BlockUsers)
	userManagment.Get("/unblock/:id", middleware.AdminMiddleware, controllers.UnBlockUsers)
	product := admin.Group("/products")
	product.Post("/add_category", controllers.AddCategory)
	product.Post("/add_product", controllers.AddProduct)
	product.Post("/delete_product/:id", controllers.DeleteProduct)
}
