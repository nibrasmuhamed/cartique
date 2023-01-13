package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nibrasmuhamed/cartique/controllers"
	"github.com/nibrasmuhamed/cartique/middleware"
)

type AdminRouter struct {
	AR controllers.AdminController
}

func NewAdminRouter(AR controllers.AdminController) *AdminRouter {
	return &AdminRouter{AR}
}

func (ar *AdminRouter) AdminRoute(admin fiber.Router) {
	admin.Post("/register", ar.AR.RegisterAdmin)
	admin.Post("/login", ar.AR.LoginAdmin)
	admin.Get("/refresh", ar.AR.RefressTokenAdmin)
	admin.Get("/logout", ar.AR.Logout)
	order := admin.Group("order")
	order.Get("view_orders", ar.AR.ViewOrders)

	userManagment := admin.Group("user_managment")
	userManagment.Get("/", middleware.AdminMiddleware, ar.AR.ViewUsers)
	userManagment.Get("/block/:id", middleware.AdminMiddleware, ar.AR.BlockUsers)
	userManagment.Get("/edit/:id", middleware.AdminMiddleware, ar.AR.EditUser)
	userManagment.Get("/unblock/:id", middleware.AdminMiddleware, ar.AR.UnBlockUsers)
	product := admin.Group("/products")
	product.Post("/add_category", ar.AR.AddCategory)
	product.Post("/add_product", ar.AR.AddProduct)
	product.Post("/delete_product/:id", ar.AR.DeleteProduct)
	product.Post("/edit_product/:id", ar.AR.EditProduct)

}
