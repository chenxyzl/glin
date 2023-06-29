package common_service

import (
	"fmt"
	"testing"
)

func TestReg(t *testing.T) {
	inputString := "@[abc#123,de1f#12387132698,哈1哈哈#7314698734] 大#1家好呀"
	pushChatMsg, purContent, numbers := ParseChatMsgFormat(inputString), ParseChatMsgPure(inputString), ParseChatMsgAt(inputString)
	fmt.Println(pushChatMsg, purContent, numbers)
}
