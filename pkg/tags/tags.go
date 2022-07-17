package tags

import (
	"fmt"
	"reflect"
)

type DecodeResult struct {
}

func Decode(in interface{}) {
	v := reflect.ValueOf(in)
	t := reflect.TypeOf(in)

	fmt.Printf("value: %+v\n", v)
	fmt.Printf("type: %+v\n", t)

	if v.Kind() == reflect.Pointer {

	}

	fmt.Println("kind", v.Elem().Kind().String())
}
