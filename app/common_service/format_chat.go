package common_service

import "strings"

func FixLaiyaChatContent(msg string) string {
	//半空格替换为全空格替换
	msg = strings.ReplaceAll(msg, "\u2005", " ")
	//去掉左右空格
	msg = strings.TrimLeft(msg, " ")
	msg = strings.TrimRight(msg, " ")
	return msg
}

func FixWechatContent(name, msg string) string {
	var atFlag = "@" + name
	if strings.Contains(msg, "\u2005") {
		atFlag += "\u2005"
	}

	idx := strings.Index(msg, atFlag)
	if idx == 0 { // >= 改为==0 //
		idx = idx + len(atFlag)
		if idx < len(msg) {
			msg = msg[idx:]
		}
	}
	//半空格替换为全空格替换
	msg = strings.ReplaceAll(msg, "\u2005", " ")
	//去掉左右空格
	msg = strings.TrimLeft(msg, " ")
	msg = strings.TrimRight(msg, " ")
	return msg
}
