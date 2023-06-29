package share

import (
	"reflect"
)

func GetTypeName[T any](v T) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	}
	return t.Name()
}
