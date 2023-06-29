package common_service

import (
	"laiya/config"
	"laiya/proto/code"
	"laiya/proto/outer"
)

var groupInfoCheck = map[outer.UpdateGroupInfo_Type]func(string) code.Code{
	outer.UpdateGroupInfo_Name: CheckGroupName,
	outer.UpdateGroupInfo_Des:  CheckGroupDes,
}

func CheckGroupInfo(typ outer.UpdateGroupInfo_Type, info string) code.Code {
	fun, ok := groupInfoCheck[typ]
	if !ok {
		return code.Code_Ok
	}
	return fun(info)
}
func CheckGroupName(name string) code.Code {
	l := len([]rune(name))
	if l < config.Get().ConstConfig.GroupNameLenMin || l > config.Get().ConstConfig.GroupNameLenMax {
		return code.Code_GroupNameLenIllegal
	}
	return code.Code_Ok
}

func CheckGroupDes(name string) code.Code {
	l := len([]rune(name))
	if l > config.Get().ConstConfig.GroupDesLenMax {
		return code.Code_GroupDesLenIllegal
	}
	return code.Code_Ok
}

func CheckVoiceRoomName(name string) code.Code {
	l := len([]rune(name))
	if l < config.Get().ConstConfig.VoiceRoomNameLenMin || l > config.Get().ConstConfig.VoiceRoomNameLenMax {
		return code.Code_GroupNameLenIllegal
	}
	return code.Code_Ok
}
