package controllers

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TestProductDB_ShowProducts(t *testing.T) {
	type fields struct {
		DB *gorm.DB
	}
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Pd := &ProductDB{
				DB: tt.fields.DB,
			}
			if err := Pd.ShowProducts(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("ProductDB.ShowProducts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
