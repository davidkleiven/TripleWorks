package pkg

import (
	"reflect"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
)

func init() {
	faker.AddProvider("uuid", func(v reflect.Value) (any, error) {
		return uuid.New(), nil
	})
}
