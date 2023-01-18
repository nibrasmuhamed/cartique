package routes

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nibrasmuhamed/cartique/controllers"
)

func TestUserRouter_Routes(t *testing.T) {
	type fields struct {
		UserRoute controllers.UserController
	}
	type args struct {
		user fiber.Router
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserRouter{
				UserRoute: tt.fields.UserRoute,
			}
			u.Routes(tt.args.user)
		})
	}
}
