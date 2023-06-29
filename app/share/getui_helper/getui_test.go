package getui_helper

import (
	"fmt"
	"testing"
)

func TestSlice(t *testing.T) {
	a := []int{1}
	b := a[:1]
	c := a[1:]
	fmt.Println(a, b, c)
}
