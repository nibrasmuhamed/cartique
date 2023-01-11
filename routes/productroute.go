package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nibrasmuhamed/cartique/controllers"
)

type ProductRouter struct {
	pr *controllers.ProductDB
}

func NewProductRoter(p *controllers.ProductDB) *ProductRouter {
	return &ProductRouter{p}
}

func (pr *ProductRouter) Routes(p fiber.Router) {
	p.Get("/", pr.pr.ShowProducts)
}
