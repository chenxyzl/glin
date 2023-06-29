package common_service

import (
	"laiya/config"
	"laiya/proto/code"
	"laiya/proto/outer"
)

var userInfoCheck = map[outer.UpdateUserInfo_Type]func(string) code.Code{
	outer.UpdateUserInfo_Name: CheckUserName,
	outer.UpdateUserInfo_Des:  CheckUserDes,
}

func CheckUserInfo(typ outer.UpdateUserInfo_Type, info string) code.Code {
	fun, ok := userInfoCheck[typ]
	if !ok {
		return code.Code_Ok
	}
	return fun(info)
}

func CheckUserName(name string) code.Code {
	l := len([]rune(name))
	if l < config.Get().ConstConfig.UserNameLenMin || l > config.Get().ConstConfig.UserNameLenMax {
		return code.Code_UserNameLenIllegal
	}
	return code.Code_Ok
}

func CheckUserDes(name string) code.Code {
	l := len([]rune(name))
	if l > config.Get().ConstConfig.UserDesLenMax {
		return code.Code_UserDesLenIllegal
	}
	return code.Code_Ok
}
