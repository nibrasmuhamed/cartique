package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nibrasmuhamed/cartique/controllers"
	"github.com/nibrasmuhamed/cartique/middleware"
)

type UserRouter struct {
	UserRoute controllers.UserController
}

func NewUserRouter(uc *controllers.UserController) *UserRouter {
	return &UserRouter{*uc}
}

func (u *UserRouter) Routes(user fiber.Router) {
	user.Post("/register", controllers.RegisterUser)
	user.Post("/login", u.UserRoute.LoginUser)
	user.Get("/verify", middleware.UserMiddleware, controllers.VerifyUser)
	user.Get("/verify/:id", middleware.UserMiddleware, controllers.VerifyUserOtp)
	user.Get("/refresh", controllers.RefreshToken)
	user.Get("/logout", controllers.Logout)
	user.Get("/", controllers.ShowProducts)

	// return app
}
