package share

import "reflect"

func GetZero[T any]() T {
	var result T
	return result
}

func IsNil(v interface{}) bool {
	rv := reflect.ValueOf(v)
	return !rv.IsValid() || rv.IsNil()
}
