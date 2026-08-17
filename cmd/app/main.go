package main

import (
	"fmt"
	"reflect"
	"strings"
)

func GetJsonTag(field reflect.StructField) {
	b, a, found := strings.Cut(field.Tag.Get("json"), ",")
	fmt.Printf("Field: %-12s | Before comma (b): %-12s | After comma (a): %-10s | Found comma: %v\n", field.Name, b, a, found)
}

type CreateBookingPayload struct {
	PropertyID *string  `validate:"required"`
	TotalPrice *float64 `validate:"required,gt=0"`
	CheckIn    *string  `validate:"required,datetime=2006-01-02"`
	CheckOut   *string  `validate:"required,datetime=2006-01-02,gtfield=CheckIn"`
}

func main() {
	fmt.Println("Hello, virtual file server")

	t := reflect.TypeOf(CreateBookingPayload{})

	// 2. Loop through every field in the struct
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		GetJsonTag(field)
	}

}
