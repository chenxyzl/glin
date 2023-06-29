package test

import (
	"fmt"
	"strings"
	"testing"
)

func TestInvited(t *testing.T) {
	text := []string{"\"小王子\"通过扫描\"运营-你得支愣起来啊\"分享的二维码加入群聊", "\"小笨魚 ꒦ິ^꒦ິ\"邀请\"黄路\"加入了群聊"}
	for _, str := range text {
		{
			tag1 := "通过扫描"
			tag2 := "分享的二维码加入群聊"
			idx1 := strings.Index(str, tag1)
			idx2 := strings.Index(str, tag2)
			if idx1 >= 0 && idx2 >= 0 && idx2 > idx1 {
				p := str[:idx1]
				op := str[idx1+len(tag1) : idx2]
				p = strings.TrimLeft(p, "\"")
				p = strings.TrimRight(p, "\"")
				op = strings.TrimLeft(op, "\"")
				op = strings.TrimRight(op, "\"")

				fmt.Println(op, p)
			}
		}
		{
			tag1 := "邀请"
			tag2 := "加入了群聊"
			idx1 := strings.Index(str, tag1)
			idx2 := strings.Index(str, tag2)
			if idx1 >= 0 && idx2 >= 0 && idx2 > idx1 {
				op := str[:idx1]
				p := str[idx1+len(tag1) : idx2]
				p = strings.TrimLeft(p, "\"")
				p = strings.TrimRight(p, "\"")
				op = strings.TrimLeft(op, "\"")
				op = strings.TrimRight(op, "\"")
				fmt.Println(op, p)
			}
		}
	}
}
