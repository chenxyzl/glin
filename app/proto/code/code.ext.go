package code

var codeMapping map[Code]string

// SetCodeMappingMsg 设置错误码对应的错误
func SetCodeMappingMsg(mapping map[Code]string) {
	if mapping != nil {
		codeMapping = mapping
	}
}

// GetCodeMsg 获取对应的错误信息
func (x Code) GetCodeMsg() string {
	if x == Code_Ok || codeMapping == nil {
		return ""
	}
	//
	if msg, ok := codeMapping[x]; ok {
		return msg
	}
	//
	return "错误码,未映射"
}
