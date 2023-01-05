package util_test

import (
	"testing"

	"github.com/nibrasmuhamed/cartique/util"
)

func TestPassword(t *testing.T) {
	got := util.Password("hellor")
	if got {
		t.Errorf("Abs(-1) = %v; want 1", got)
	}
	got = util.Password("Hello@123")
	if !got {
		t.Error("this should pass")
	}

}
