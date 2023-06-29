package code

import (
	"fmt"
	"google.golang.org/protobuf/reflect/protoreflect"
	"testing"
)

func TestCode(t *testing.T) {
	v := Code_Ok
	var b interface{} = v
	v1, o := b.(protoreflect.Enum)
	if !o {
		t.Error()
	} else {
		fmt.Println(v1.Number())
	}
}
