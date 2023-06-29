package component

import (
	"fmt"
	"laiya/common_service"
	"strconv"
	"strings"
	"testing"
)

func TestContent(t *testing.T) {
	content := "@[二鸭-厨房摇人#300002] 11"
	tag := "#" + strconv.Itoa(int(300002)) + "]"
	idx := strings.Index(content, tag)
	if idx < 0 {
		t.Error()
	}
	// 012ab56
	realContent := string(content[idx+len(tag):])
	fmt.Println(realContent)
	realContent = common_service.FixLaiyaChatContent(realContent)
	fmt.Println(realContent)
}
