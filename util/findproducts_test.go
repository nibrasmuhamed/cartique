package util

import (
	"reflect"
	"testing"

	"github.com/nibrasmuhamed/cartique/models"
)

func TestFindProducts(t *testing.T) {
	type args struct {
		a []int
	}
	tests := []struct {
		name string
		args args
		want []models.ProductRespHome
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindProducts(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindProducts() = %v, want %v", got, tt.want)
			}
		})
	}
}
