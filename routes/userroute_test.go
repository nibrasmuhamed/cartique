package routes

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRoute(t *testing.T) {
	// Create a new Fiber app
	app := fiber.New()
	u := UserRouter{}
	u.Routes(app)

	// Register the route

	// Create a request to send to the app
	req := httptest.NewRequest("GET", "/user", nil)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to record the response
	respRecorder := httptest.NewRecorder()

	// Serve the request
	app.Test(req, 1)

	// Check the status code
	if respRecorder.Code != http.StatusOK {
		t.Errorf("Expected status code to be %d, but got %d", http.StatusOK, respRecorder.Code)
	}
	fmt.Println("  hai  ", respRecorder.Body.String())
	// Check the response body
	expected := "user not authorized"
	if respRecorder.Body.String() != expected {
		t.Errorf("Expected response body to be %q, but got %q", expected, respRecorder.Body.String())
	}
}
