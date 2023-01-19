package routes_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nibrasmuhamed/cartique/controllers"
	"github.com/nibrasmuhamed/cartique/database"
	"github.com/nibrasmuhamed/cartique/routes"
	"github.com/stretchr/testify/assert"
)

func TestUserRouter_Routes(t *testing.T) {
	tests := []struct {
		description string

		// Test input
		route string

		// Expected output
		expectedError bool
		expectedCode  int
		expectedBody  string
	}{
		{
			description:   "index route",
			route:         "/user/login",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  "OK",
		},
		{
			description:   "non existing route",
			route:         "/i-dont-exist",
			expectedError: false,
			expectedCode:  404,
			expectedBody:  "Cannot GET /i-dont-exist",
		},
	}

	var UC *controllers.UserController
	var AC *controllers.AdminController
	var ar *routes.AdminRouter
	var ur *routes.UserRouter
	var Pc *controllers.ProductDB
	var Pr *routes.ProductRouter

	//  := routes.UserRouter
	app := fiber.New()
	x := database.InitDB()
	AC = controllers.NewAdminController(x)
	UC = controllers.NewUserController(x)
	Pc = controllers.NewProductDB(x)
	ur = routes.NewUserRouter(UC)
	ar = routes.NewAdminRouter(*AC)
	Pr = routes.NewProductRoter(Pc)

	products := app.Group("/products")
	Pr.Routes(products)
	admin := app.Group("/admin")
	ar.AdminRoute(admin)
	user := app.Group("/user")
	ur.Routes(user)

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route
		// from the test case
		req, _ := http.NewRequest(
			"GET",
			test.route,
			nil,
		)

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		// verify that no error occured, that is not expected

		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next
		// test case needs to be processed
		// if test.expectedError {
		// 	continue
		// }

		// Verify if the status code is as expected
		// assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		// assert.Nilf(t, err, test.description)

		// Verify, that the reponse body equals the expected body
		assert.Equalf(t, test.expectedBody, string(body), test.description)
	}
}
