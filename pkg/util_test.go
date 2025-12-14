package pkg

import (
	"errors"
	"testing"
)

func TestMustPanicOnErr(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Wanted panic")
		}
	}()

	Must(0, errors.New("Something went wrong"))
}
