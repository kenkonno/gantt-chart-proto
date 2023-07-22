package utils

import (
	"fmt"
	"github.com/golang-module/carbon/v2"
	"testing"
)

func TestStrToDate(t *testing.T) {
	got := StrToDate("2023-01-23")
	want := carbon.NewCarbon().SetDate(2023, 1, 1)
	if got != want {
		t.Error(fmt.Sprintf("%v, %v", want, got))
	}

}
