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
	user.Post("/register", u.UserRoute.RegisterUser)
	user.Post("/login", u.UserRoute.LoginUser)
	user.Get("/verify", middleware.UserMiddleware, u.UserRoute.VerifyUser)
	user.Get("/verify/:id", middleware.UserMiddleware, u.UserRoute.VerifyUserOtp)
	user.Get("/refresh", u.UserRoute.RefreshToken)
	user.Get("/logout", u.UserRoute.Logout)
	user.Get("/", u.UserRoute.ShowProducts)
	user.Get("/addtocart/:id", u.UserRoute.AddToCart)
	user.Get("/showcart", u.UserRoute.ShowCart)
	user.Post("/add_address", middleware.UserMiddleware, u.UserRoute.AddAddress)
	user.Get("/show_address", middleware.UserMiddleware, u.UserRoute.ShowAddress)

}
