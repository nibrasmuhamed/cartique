package util

import (
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestInternalServerErr(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InternalServerErr(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("InternalServerErr() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
