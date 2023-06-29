package common_service

import (
	"regexp"
	"strconv"
	"strings"
)

// 匹配@[~]之间的的所有字符串
var reInFrame = regexp.MustCompile(`@\[(.*?)\]`) //
// 匹配#号后面的Id
var reId = regexp.MustCompile(`#.*$`) //
// 匹配@[~]之间的的所有字符串
var reAt = regexp.MustCompile(`@\[([^]]+)\]`) //
// 匹配@[~]在内的所有字符串
var reWithFrame = regexp.MustCompile(`@\[.*?\]`) //

func ParseChatMsgFormat(str string) string {
	// 正则表达式匹配 @[...] 中的内容
	matches := reInFrame.FindStringSubmatch(str)
	if len(matches) < 2 {
		return str
	}
	// 获取 @[...] 中的内容
	content := matches[1]
	// 去除 @[ 和 ] 符号
	content = strings.TrimPrefix(content, "@[")
	content = strings.TrimSuffix(content, "]")
	// 以逗号分割子串
	substrs := strings.Split(content, ",")
	// 处理每个子串
	for i, substr := range substrs {
		// 去除 # 及之后的数字
		substrs[i] = reId.ReplaceAllString(substr, "")
		// 在每个子串前面加上 @ 符号
		substrs[i] = "@" + substrs[i]
	}
	// 将所有子串用空格连接为一个字符串
	result := strings.Join(substrs, "")
	// 将处理后的字符串替换原始字符串中的 @[...] 部分
	result = strings.Replace(str, "@["+content+"]", result, 1)
	return result
}

// ParseChatMsgAt 解析At的对象
func ParseChatMsgAt(context string) []uint64 {
	match := reAt.FindStringSubmatch(context)
	if len(match) < 2 {
		return nil
	}
	//
	numbers := make([]uint64, 0)
	strList := strings.Split(match[1], ",")
	for _, str := range strList {
		if num, err := strconv.ParseUint(strings.TrimRight(strings.Split(str, "#")[1], "1"), 10, 64); err == nil {
			numbers = append(numbers, num)
		}
	}
	//
	return numbers
}

// ParseChatMsgPure 解析纯净的聊天内容
func ParseChatMsgPure(context string) string {
	//标准化空格,前后去掉多余空格
	context = FixLaiyaChatContent(context)
	//干净的聊天内容
	purContent := reWithFrame.ReplaceAllString(context, "")
	purContent = strings.TrimPrefix(purContent, " ")
	purContent = strings.TrimSuffix(purContent, " ")

	return purContent
}
